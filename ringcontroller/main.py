#!/usr/bin/env python3
import sys
import base64
import io
import logging

from os import path
from swift.common.ring import RingBuilder
from swift.common.ring.utils import parse_add_value
from kubernetes import client, config
from kubernetes.config import ConfigException
from kubernetes.client.rest import ApiException
from ring_controller import RingController

logger = logging.getLogger('ringcontroller')


def load_config_map(v1, namespace_name):
    config_map = None
    l = namespace_name.split('/')
    logger.info("loading config_map %s/%s", l[0], l[1])
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


def patch_config_map(v1, config_map_meta, ring_data):
    logger.info("patching config_map %s/%s", config_map_meta.namespace,
                config_map_meta.name)
    v1.patch_namespaced_config_map(
        namespace=config_map_meta.namespace, name=config_map_meta.name, body=client.V1ConfigMap(
            binary_data=ring_data,
        ))


def setup_logging():
    logger.setLevel(logging.DEBUG)
    handler = logging.StreamHandler(sys.stdout)
    formatter = logging.Formatter(
        '%(asctime)s - %(funcName)s() - %(levelname)s - %(message)s')
    handler.setFormatter(formatter)
    logger.addHandler(handler)


def main(argv):
    setup_logging()

    not_enough_arguments = len(argv) < 4
    if not_enough_arguments:
        print("use reconcile-ring.py <config_map> <type> <devices>")
        print("for example: reconcile-ring.py contrail/swift-object object r1z1-192.168.0.2:6000/d3 r1z2-192.168.2.2:5000/d1")
        print("each device has format rREGIONzZONE-IP:PORT/DEVICE")
        exit(1)

    logger.debug(argv)

    load_config()
    v1 = client.CoreV1Api()
    config_map_namespace_name = argv[1]
    ring_type = argv[2]
    device_strings = argv[3:]
    config_map = load_config_map(v1, config_map_namespace_name)

    if config_map.binary_data != None and ring_type in config_map.binary_data:
        builder_ = base64.b64decode(config_map.binary_data[ring_type])
        logger.info("loading existing ring builder")
        builder = RingBuilder.load("", open=lambda a, b: io.BytesIO(builder_))
    else:
        logger.info("creating a new ring builder")
        builder = RingBuilder(10, 1, 1)

    r = RingController(builder, ring_type, logger=logger)
    logger.info("reconciling ring")
    r.reconcile(device_strings)
    ring_data = r.get_ring_data()

    patch_config_map(v1, config_map.metadata, ring_data)


if __name__ == "__main__":
    # execute only if run as a script
    main(sys.argv)
