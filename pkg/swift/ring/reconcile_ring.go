package ring

var reconcileRingScript = `#!/var/lib/kolla/venv/bin/python
from sys import argv
from os import path

from swift.common.ring import RingBuilder
from swift.common.ring.utils import parse_add_value


def format_device(d):
    return 'r%(region)sz%(zone)s-%(ip)s:%(port)s/%(device)s' % d


not_enough_arguments = len(argv) < 3
if not_enough_arguments:
    print("use reconcile-ring.py <builder_file> <devices>")
    print("for example: reconcile-ring.py swift.ring.builder r1z1-192.168.0.2:6000/d3 r1z2-192.168.2.2:5000/d1")
    print("each device has format rREGIONzZONE-IP:PORT/DEVICE")
    exit(1)

device_strings = argv[2:]

builder_filename = argv[1]
if not builder_filename.endswith('.builder'):
    ring_filename = builder_filename
else:
    ring_filename = builder_filename[:-len('.builder')]
ring_filename += '.ring.gz'

if path.exists(builder_filename):
    print("opening file " + builder_filename)
    builder = RingBuilder.load(builder_filename)
else:
    print("creating file " + builder_filename)
    # FIXME: Make replicas count configurable
    builder = RingBuilder(10, 1, 1)

# Add new devices
for dev_string in device_strings:
    dev = parse_add_value(dev_string)
    if len(builder.search_devs(dev)) == 0:
        dev['weight'] = 1
        dev_id = builder.add_dev(dev)
        print("Adding " + dev_string)

# Remove devices
for dev in builder.search_devs([]):
    dev_string = format_device(dev)
    if dev_string not in device_strings:
        builder.remove_dev(dev['id'])
        print("Removing " + dev_string)

builder.rebalance()
builder.save(builder_filename)
builder.get_ring().save(ring_filename)
`
