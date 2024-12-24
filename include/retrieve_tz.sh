#!/usr/bin/env sh
{ test -z "$REFRESH_TZ" && test -f /data/timezonedb.csv; } && exit 0
if ! test -d /data
then mkdir /data || exit 1
fi
cd /tmp || exit 1
curl -o tzdb.csv.zip https://timezonedb.com/files/TimeZoneDB.csv.zip &&
unzip tzdb.csv.zip &&
echo 'id,zone_name,country_code,abbreviation,time_start,gmt_offset,dst' > /data/timezonedb.csv; \
{
  cat -b time_zone.csv |
    tr '\t' ',' |
    sed -E 's/^[ \t]+//';
} >> /data/timezonedb.csv

