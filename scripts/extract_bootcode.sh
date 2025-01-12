dd bs=1 skip=64 count=4032 if=$1 of=$2
echo $(md5 $2)
