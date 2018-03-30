#!/bin/sh

# Usage: stresster_tv.sh test.jmx threads rampup counter


if [ "$#" -ne 4 ]; then
  echo "Usage: $0 test.jmx threads rampup counter"
  exit 1
fi

AWS1=yourinstance1.compute.amazonaws.com
AWS2=yourinstance2.compute.amazonaws.com
AWS3=yourinstance3.compute.amazonaws.com

USER=ec2-user
KEY=key.pem
HEAP="-Xms512m -Xmx2048m -Xss228k"
JMETER=apache-jmeter-3.1/bin/ApacheJMeter.jar

echo "Distributing the test plans..."
scp -i $KEY $1 $USER@$AWS1:
scp -i $KEY $1 $USER@$AWS2:
scp -i $KEY $1 $USER@$AWS3:

echo " "
echo "Deleting existing reports..."
xterm -e ssh -i $KEY $USER@$AWS1 -t "command; rm aws1.jtl" &
xterm -e ssh -i $KEY $USER@$AWS2 -t "command; rm aws2.jtl" &
xterm -e ssh -i $KEY $USER@$AWS3 -t "command; rm aws3.jtl" &

COUNT=$(($2 / 3))
COUNTER1=$4
COUNTER2=$(((1 * $COUNT) + $4))
COUNTER3=$(((2 * $COUNT) + $4))

if [ "$COUNT" -gt 0 ]; then
  echo "Starting the tests, $COUNT threads per node"
  xterm -e ssh -i $KEY $USER@$AWS1 -t "command; java $HEAP -jar $JMETER -n -l test_part1.jtl -t $1 -Jthreads=$COUNT -Jrampup=$3 -Jcounter=$COUNTER1" &
  xterm -e ssh -i $KEY $USER@$AWS2 -t "command; java $HEAP -jar $JMETER -n -l test_part2.jtl -t $1 -Jthreads=$COUNT -Jrampup=$3 -Jcounter=$COUNTER2" &
  xterm -e ssh -i $KEY $USER@$AWS3 -t "command; java $HEAP -jar $JMETER -n -l test_part3.jtl -t $1 -Jthreads=$COUNT -Jrampup=$3 -Jcounter=$COUNTER3" &
else
  echo "Distribution only, no tests started"
fi
