#!/bin/bash
source /etc/profile
#usage ./port_status.sh localhost 80 8080 8088
#parameter 1 is host
#parameters after 1 are list of ports to check
#all ports are ok echo 1, else echo 0 

function clean_files(){
  dumpfile_nums=$((`ls resourcemanager_fullgc*|wc -l` -2))
  if [ $dumpfile_nums -gt 2 ]; then
    ls -tr resourcemanager_fullgc*|head -${dumpfile_nums}|xargs rm -f
  fi
}

if [ $# -lt 2 ] ; then
  echo "USAGE: $0 host [ports]" 
  exit
fi
 
address=$1
ret=0
shift
i=$@


for i in "$@"; do
  if command -v nc >/dev/null 2>&1; then
 # echo "exists nc"
   #  echo  $address $i
     nc -w 10  $address $i  < /dev/null >/dev/null 2>&1 
     
   #  echo status ======= $?
     if [ $? -eq 0 ] ; then
        ret=$(( $ret + 0 ))
     else
        ret=$(( $ret + 1 ))
     fi
  else
    ret=2
  fi
done


now_date=`date +%Y-%m-%d#%H:%M:%S`
pid=`jps | grep ResourceManager | awk '{print $1}'`
fgc=`$JAVA_HOME/bin/jstat -gcutil $pid |grep -v FGC|awk '{print $9}'`
fgct=`$JAVA_HOME/bin/jstat -gcutil $pid|grep -v FGC|awk '{print $10}'`
oldgc=`$JAVA_HOME/bin/jstat -gcutil $pid|grep -v FGC|awk '{print $4}'`


gc_nums=`sed '/^NUMS=/!d;s/.*=//' gcconfig`
r_fgc=`sed '/^FGC=/!d;s/.*=//' gcconfig`
r_fgct=`sed '/^FGCT=/!d;s/.*=//' gcconfig`
dumps=`sed '/^DUMPS=/!d;s/.*=//' gcconfig`
oldgc_nums=`sed '/^OLDGCNUMS=/!d;s/.*=//' gcconfig`


sed -i 's#^FGC=\(.*\)#FGC='"${fgc}"'#g' gcconfig
sed -i 's#^FGCT=\(.*\)#FGCT='"${fgct}"'#g' gcconfig

if [ $dumps -gt 0 -a $dumps -le 20 ]; then
   ((dumps=$dumps+1))
   sed -i 's#^DUMPS=[0-9]*#DUMPS='"${dumps}"'#g' gcconfig
   exit 1
elif [ $dumps -gt 20 ]; then
   sed -i 's#^DUMPS=[0-9]*#DUMPS=0#g' gcconfig
   kill -9 ${pid}
fi

if [ $gc_nums -gt 10 ]; then
   clean_files
   sed -i 's#^DUMPS=[0-9]*#DUMPS=1#g' gcconfig
   sed -i 's#^NUMS=[0-9]*#NUMS=0#g' gcconfig
   $JAVA_HOME/bin/jmap -dump:format=b,file=resourcemanager_fullgc_${pid}_${now_date} $pid
   sed -i 's#^DUMPS=[0-9]*#DUMPS=0#g' gcconfig
   kill -9 $pid
elif [ $fgc -gt $r_fgc ]; then
  ((gc_nums=$gc_nums+1))
  ret=$(( $ret + 1 ))
  sed -i 's#^NUMS=[0-9]*#NUMS='"${gc_nums}"'#g' gcconfig
else
  sed -i 's#^NUMS=[0-9]*#NUMS=0#g' gcconfig
fi

if [ $oldgc_nums -gt 10 ]; then
  sed -i 's#^DUMPS=[0-9]*#DUMPS=1#g' gcconfig
  sed -i 's#^NUMS=[0-9]*#NUMS=0#g' gcconfig
  sed -i 's#^OLDGCNUMS=[0-9]*#OLDGCNUMS=0#g' gcconfig
  kill -9 $pid
elif [ $(echo "$oldgc>96"| bc) = 1 ]; then
  ((OLDGCNUMS=$OLDGCNUMS+1))
  sed -i 's#^OLDGCNUMS=[0-9]*#OLDGCNUMS='"${OLDGCNUMS}"'#g' gcconfig
else
  sed -i 's#^OLDGCNUMS=[0-9]*#OLDGCNUMS=0#g' gcconfig
fi

exit $ret

