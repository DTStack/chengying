#!/bin/bash

source /etc/profile
ret=0
bin="$HADOOP_HOME"/bin/yarn
SPARK_REMARK='/tmp/spark.remark'

function clean_files(){
  dumpfile_nums=$((`ls thriftserver_fullgc*|wc -l` -2))
  if [ $dumpfile_nums -gt 2 ]; then
    ls -tr thriftserver_fullgc*|head -${dumpfile_nums}|xargs rm -f
  fi
}

function run(){
    remark=$(GetRemark ${SPARK_REMARK})
    num=$("$bin" application -list|grep ThriftServer02|grep -v grep|grep RUNNING|wc -l)
    if [ $num -ne 1 ]
    then
      ret=$(( $ret + 1 ))
      if [ "$remark" = "" ] ;then
        echo 1 > "$SPARK_REMARK"
      elif [ $remark -ge 3 ]; then
        ps -ef |grep start_spark|grep -v grep|cut -c 9-15|xargs kill -9
        echo 0 > "$SPARK_REMARK"
      else
        remark=$(( $remark+1 ))
        echo $remark > "$SPARK_REMARK"
      fi
    else
      echo 0 > "$SPARK_REMARK"
    fi
}

function GetRemark() {
    local remark
    if [ -f "$1" ] && [ -s "$1" ]; then
        remark=$(cat $1)

        if [ $remark -gt 3 ]; then
            rm -f $1
            remark=""
        fi
    else
        remark=""
    fi
    echo $remark
}
run

now_date=`date +%Y-%m-%d#%H:%M:%S`
pid=`ps -ef |grep thriftserver|grep -v grep|grep -v "start_spark"|awk '{print $2}'`
fgc=`$JAVA_HOME/bin/jstat -gcutil $pid |grep -v FGC|awk '{print $9}'`
fgct=`$JAVA_HOME/bin/jstat -gcutil $pid|grep -v FGC|awk '{print $10}'`


gc_nums=`sed '/^NUMS=/!d;s/.*=//' gcconfig`
r_fgc=`sed '/^FGC=/!d;s/.*=//' gcconfig`
r_fgct=`sed '/^FGCT=/!d;s/.*=//' gcconfig`
dumps=`sed '/^DUMPS=/!d;s/.*=//' gcconfig`

sudo sed -i 's#^FGC=[0-9]*#FGC='"${fgc}"'#g' gcconfig
sudo sed -i 's#^FGCT=*.*#FGCT='"${fgct}"'#g' gcconfig

if [ $dumps -gt 0 -a $dumps -le 20 ]; then
   ((dumps=$dumps+1))
   sudo sed -i 's#^DUMPS=[0-9]*#DUMPS='"${dumps}"'#g' gcconfig
   exit 1
elif [ $dumps -gt 20 ]; then
   sudo sed -i 's#^DUMPS=[0-9]*#DUMPS=0#g' gcconfig
   sudo kill -9 ${pid}
fi

if [ $gc_nums -gt 10 ]; then
   clean_files
   sudo sed -i 's#^DUMPS=[0-9]*#DUMPS=1#g' gcconfig
   sudo sed -i 's#^NUMS=[0-9]*#NUMS=0#g' gcconfig
   $JAVA_HOME/bin/jmap -dump:format=b,file=thriftserver_fullgc_${pid}_${now_date} $pid
   sudo sed -i 's#^DUMPS=[0-9]*#DUMPS=0#g' gcconfig
   sudo kill -9 $pid
elif [ $fgc -gt $r_fgc ]; then
  ((gc_nums=$gc_nums+1))
  ret=$(( $ret + 1 ))
  sudo sed -i 's#^NUMS=[0-9]*#NUMS='"${gc_nums}"'#g' gcconfig
else
  sudo sed -i 's#^NUMS=[0-9]*#NUMS=0#g' gcconfig
fi
exit $ret

