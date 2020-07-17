#!/usr/bin/env python3
import argparse
import base64
import io
import logging
import sys

from kubernetes import client, config
from kubernetes.config import ConfigException
from swift.common.ring import RingBuilder

from ring_controller import RingController

logger = logging.getLogger('ringcontroller')


def load_config_map(v1, namespace_name):
    ns, name = namespace_name.split('/')
    logger.info("loading config_map %s/%s", ns, name)
    return v1.read_namespaced_config_map(namespace=ns, name=name)


def load_config():
    try:
        config.load_incluster_config()
    except ConfigException:
        config.load_kube_config()


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


def read_args():
    parser = argparse.ArgumentParser(
        description="It reconciles swift ring that is stored in k8s config map",
        formatter_class=argparse.RawDescriptionHelpFormatter,
        epilog='''
example usage:
  main.py contrail/swift-object object r1z1-192.168.0.2:6000/d3 r1z2-192.168.2.2:5000/d1'''
    )
    parser.add_argument("config_map_name", help="config map namespace/name")
    parser.add_argument("ring_type", help="ring type")
    parser.add_argument(
        "devices", nargs='+', help="list of devices in format: rREGIONzZONE-IP:PORT/DEVICE")

    return parser.parse_args()


def main(argv):
    setup_logging()
    args = read_args()

    logger.debug(args)
    load_config()
    v1 = client.CoreV1Api()
    config_map = load_config_map(v1, args.config_map_name)

    if config_map.binary_data != None and args.ring_type in config_map.binary_data:
        builder_ = base64.b64decode(config_map.binary_data[args.ring_type])
        logger.info("loading existing ring builder")
        builder = RingBuilder.load("", open=lambda a, b: io.BytesIO(builder_))
    else:
        logger.info("creating a new ring builder")
        builder = RingBuilder(10, 1, 1)

    r = RingController(builder, args.ring_type, logger=logger)
    logger.info("reconciling ring")
    r.reconcile(args.devices)
    ring_data = r.get_ring_data()

    patch_config_map(v1, config_map.metadata, ring_data)


if __name__ == "__main__":
    main(sys.argv)
