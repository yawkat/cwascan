#!/usr/bin/python3

import psycopg2
import sys
import PIL.Image
import datetime
import colorsys

SECONDS_PER_PIXEL = 10

conn = psycopg2.connect(sys.argv[1])
cur = conn.cursor()

colors = {}
next_hue = 0.0

cur.execute("select distinct scannerid from rpik")
for scanner_id, in cur.fetchall():
    cur.execute("select min(timestamp), max(timestamp) from rpik where scannerid = %s", (scanner_id,))
    lo: datetime.datetime
    hi: datetime.datetime
    (lo, hi) = cur.fetchone()
    image = PIL.Image.new("RGB", (int((hi.timestamp() - lo.timestamp()) / SECONDS_PER_PIXEL) + 1, 110), color="black")
    pixel_map = image.load()
    cur.execute("select timestamp, rpik, rssi from rpik where scannerid = %s order by timestamp", (scanner_id,))
    while True:
        row = cur.fetchone()
        if row is None:
            break
        (ts, rpik, rssi) = row
        x = int((ts.timestamp() - lo.timestamp()) / SECONDS_PER_PIXEL)
        y = int(-rssi)
        if rpik not in colors:
            (r, g, b) = colorsys.hsv_to_rgb(next_hue, 1, 1)
            colors[rpik] = (int(r * 255), int(g * 255), int(b * 255))
            next_hue = (next_hue + 0.3) % 1
        pixel_map[x, y] = colors[rpik]

    image.save("plot%s.png" % scanner_id)

