import base64
import io
import os
import tempfile

from swift.common.ring.utils import parse_add_value
from gzip import GzipFile


class RingController(object):
    def __init__(self, builder, ring_type):
        self.builder = builder
        self.ring_type = ring_type

    def reconcile(self, devices):
        def format_device(d):
            return 'r%(region)sz%(zone)s-%(ip)s:%(port)s/%(device)s' % d

         # Add new devices
        for dev_string in devices:
            dev = parse_add_value(dev_string)
            if len(self.builder.search_devs(dev)) == 0:
                dev['weight'] = 1
                self.builder.add_dev(dev)
                print("Adding " + dev_string)

        # Remove devices
        for dev in self.builder.search_devs([]):
            dev_string = format_device(dev)
            if dev_string not in devices:
                self.builder.remove_dev(dev['id'])
                print("Removing " + dev_string)

        self.builder.rebalance()

    def _serialize_ring(self, filename):
        buf = io.BytesIO()
        # Override the timestamp so that the same ring data creates
        # the same bytes on disk. This makes a checksum comparison a
        # good way to see if two rings are identical.
        with GzipFile(filename=filename, fileobj=buf, mode='wb', mtime=1300507380.0, compresslevel=9) as f:
            self.builder.get_ring().serialize_v1(f)

        return base64.b64encode(buf.getvalue()).decode("utf-8")

    def _serialize_builder(self):
        builder_ = None
        handle, filename = tempfile.mkstemp()
        os.close(handle)
        try:
            self.builder.save(filename)
            with open(filename, 'rb') as fp:
                builder_ = base64.b64encode(fp.read()).decode("utf-8")
        finally:
            os.remove(filename)
        return builder_

    def get_ring_data(self):
        ring_filename = self.ring_type + '.ring.gz'
        return {
            ring_filename: self._serialize_ring(ring_filename),
            self.ring_type: self._serialize_builder(),
        }
