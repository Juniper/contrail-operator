#!/usr/bin/env python3
import sys
import base64
import io

from os import path
from swift.common.ring import RingBuilder
from swift.common.ring.utils import parse_add_value
from kubernetes import client, config
from kubernetes.config import ConfigException
from kubernetes.client.rest import ApiException


def format_device(d):
    return 'r%(region)sz%(zone)s-%(ip)s:%(port)s/%(device)s' % d


def load_config_map(v1, namespace_name):

    config_map = None
    l = namespace_name.split('/')
    try:
        config_map = v1.read_namespaced_config_map(
            namespace=l[0], name=l[1])

    except ApiException as e:
        print("api error: %s" % e)
        exit(1)

    return config_map


def load_config():
    try:
        config.load_incluster_config()
    except ConfigException:
        config.load_kube_config()
        pass


def reconcile(builder, device_strings):
   # Add new devices
    for dev_string in device_strings:
        dev = parse_add_value(dev_string)
        if len(builder.search_devs(dev)) == 0:
            dev['weight'] = 1
            builder.add_dev(dev)
            print("Adding " + dev_string)

    # Remove devices
    for dev in builder.search_devs([]):
        dev_string = format_device(dev)
        if dev_string not in device_strings:
            builder.remove_dev(dev['id'])
            print("Removing " + dev_string)

    builder.rebalance()


def patch_config_map(v1, namespace_name, builder, builder_name):

    l = namespace_name.split('/')
    ring_filename = builder_name + '.ring.gz'
    builder.save(builder_name)
    builder.get_ring().save(ring_filename)

    with open(ring_filename, 'rb') as reader:
        ring = base64.b64encode(reader.read())

    with open(builder_name, 'rb') as reader:
        builder_file = base64.b64encode(reader.read())

    v1.patch_namespaced_config_map(
        namespace=l[0], name=l[1], body=client.V1ConfigMap(
            binary_data={
                ring_filename: ring.decode("utf-8"),
                builder_name: builder_file.decode("utf-8"),
            }
        ))


def main(argv):

    not_enough_arguments = len(argv) < 4
    if not_enough_arguments:
        print("use reconcile-ring.py <config_map> <type> <devices>")
        print("for example: reconcile-ring.py contrail/swift-object object r1z1-192.168.0.2:6000/d3 r1z2-192.168.2.2:5000/d1")
        print("each device has format rREGIONzZONE-IP:PORT/DEVICE")
        exit(1)

    load_config()
    v1 = client.CoreV1Api()
    config_map_namespace_name = argv[1]
    builder_name = argv[2]
    device_strings = argv[3:]
    config_map = load_config_map(v1, config_map_namespace_name)

    if config_map.binary_data != None and builder_name in config_map.binary_data:
        builder_ = base64.b64decode(config_map.binary_data[builder_name])
        builder = RingBuilder.load("", open=lambda a, b: io.BytesIO(builder_))
    else:
        builder = RingBuilder(10, 1, 1)

    reconcile(builder, device_strings)

    patch_config_map(v1, config_map_namespace_name, builder, builder_name)


if __name__ == "__main__":
    # execute only if run as a script
    main(sys.argv)
