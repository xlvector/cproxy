import urllib2, sys

for line in sys.stdin:
    urllib2.urlopen(line.strip()).read()