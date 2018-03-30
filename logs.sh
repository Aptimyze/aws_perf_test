#!/bin/sh

AWS1=yourinstance1.compute.amazonaws.com
AWS2=yourinstance2.compute.amazonaws.com
AWS3=yourinstance3.compute.amazonaws.com

USER=ec2-user
KEY=key.pem

DATE=$(date +"%Y%m%d%H%M")
mkdir report/$DATE

echo "Getting the test logs..."
scp -i $KEY $USER@$AWS1:aws1.jtl report/$DATE/.
scp -i $KEY $USER@$AWS2:aws2.jtl report/$DATE/.
scp -i $KEY $USER@$AWS3:aws3.jtl report/$DATE/.

echo " "
echo "Done."
