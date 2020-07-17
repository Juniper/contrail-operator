import base64
import io
import unittest
from gzip import decompress

from swift.common.ring import RingBuilder
from swift.common.ring.utils import parse_add_value

from ring_controller import RingController


class TestRingController(unittest.TestCase):

    def test_reconcile_add_devices(self):
        # given a empty ring builder
        b = RingBuilder(10, 1, 1)
        r = RingController(b, "object")

        # when new devices are added
        ds = ["r1z1-192.168.0.2:6000/d3", "r1z1-192.168.2.2:5000/d1"]
        r.reconcile(ds)

        # then new devices are present
        devs = b.search_devs([])
        self.assertEqual(2, len(devs))
        self.assertEqual(devs,
                         [{
                             'device': 'd3', 'id': 0, 'ip': '192.168.0.2', 'meta': '',
                             'parts': 512, 'parts_wanted': 0, 'port': 6000, 'region': 1,
                             'replication_ip': None, 'replication_port': None, 'weight': 1.0, 'zone': 1
                         }, {
                             'device': 'd1', 'id': 1, 'ip': '192.168.2.2', 'meta': '',
                             'parts': 512, 'parts_wanted': 0, 'port': 5000, 'region': 1,
                             'replication_ip': None, 'replication_port': None, 'weight': 1.0, 'zone': 1
                         }])

    def test_reconcile_remove_device(self):
        # given a ring builder with two devices
        b = RingBuilder(10, 1, 1)
        r = RingController(b, "object")
        ds = ["r1z1-192.168.0.2:6000/d3", "r1z2-192.168.2.2:5000/d1"]
        r.reconcile(ds)

        # when device is removed
        ds = ["r1z2-192.168.2.2:5000/d1"]
        r.reconcile(ds)

        # then only one device is present
        devs = b.search_devs([])
        self.assertEqual(1, len(devs))
        self.assertEqual(devs,
                         [{
                             'device': 'd1', 'id': 1, 'ip': '192.168.2.2', 'meta': '',
                             'parts': 1024, 'parts_wanted': 0, 'port': 5000, 'region': 1,
                             'replication_ip': None, 'replication_port': None, 'weight': 1.0, 'zone': 2
                         }])

    def test_reconcile_change_device(self):
        # given a ring builder with two devices
        b = RingBuilder(10, 1, 1)
        r = RingController(b, "object")
        ds = ["r1z1-192.168.0.2:6000/d3", "r1z2-192.168.2.2:5000/d1"]
        r.reconcile(ds)

        # when one device is changed
        ds = ["r1z1-192.168.0.3:6000/d3", "r1z2-192.168.2.2:5000/d1"]
        r.reconcile(ds)

        # then device is updated
        devs = b.search_devs([])
        self.assertEqual(2, len(devs))
        self.assertEqual(devs,
                         [{
                             'device': 'd1', 'id': 1, 'ip': '192.168.2.2', 'meta': '',
                             'parts': 512, 'parts_wanted': 0, 'port': 5000, 'region': 1,
                             'replication_ip': None, 'replication_port': None, 'weight': 1.0, 'zone': 2
                         }, {
                             'device': 'd3', 'id': 2, 'ip': '192.168.0.3', 'meta': '',
                             'parts': 512, 'parts_wanted': 0, 'port': 6000, 'region': 1,
                             'replication_ip': None, 'replication_port': None, 'weight': 1.0, 'zone': 1
                         }])

    def test_get_ring_data(self):
        # given a ring builder with two devices
        b = RingBuilder(10, 1, 1)
        r = RingController(b, "object")
        ds = ["r1z1-192.168.0.2:6000/d3", "r1z2-192.168.2.2:5000/d1"]
        r.reconcile(ds)

        # when ring is serialized
        ring_data = r.get_ring_data()

        # then ring and builder is returned
        self.assertIn("object.ring.gz", ring_data)
        self.assertIn("object", ring_data)

        # then correct compressed ring is returned
        object_ring = decompress(base64.b64decode(ring_data["object.ring.gz"]))
        expected_ring = b'{"byteorder": "little", "devs": [{"device": "d3", "id": 0, "ip": "192.168.0.2", "meta": "", "port": 6000, "region": 1, "replication_ip": null, "replication_port": null, "weight": 1.0, "zone": 1}, {"device": "d1", "id": 1, "ip": "192.168.2.2", "meta": "", "port": 5000, "region": 1, "replication_ip": null, "replication_port": null, "weight": 1.0, "zone": 2}], "part_shift": 22, "replica_count": 1, "version": 3}'
        self.assertIn(expected_ring, object_ring)

        # then correct serialized builder is returned
        builder_ = base64.b64decode(ring_data["object"])
        loadedBuilder = RingBuilder.load(
            "", open=lambda a, b: io.BytesIO(builder_))
        self.assertEqual(b.devs, loadedBuilder.devs)


if __name__ == '__main__':
    unittest.main()
