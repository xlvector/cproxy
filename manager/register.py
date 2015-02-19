import urllib2, sys

for line in sys.stdin:
    try:
        urllib2.urlopen(line.strip()).read()
    except Exception, e:
        pass
