import sys
import re
from awsglue.transforms import *
from awsglue.utils import getResolvedOptions
from pyspark.context import SparkContext
from awsglue.context import GlueContext
from awsglue.job import Job
from awsglue.dynamicframe import DynamicFrame
from datetime import datetime, timedelta

glueContext = GlueContext(SparkContext.getOrCreate())

args = getResolvedOptions(sys.argv, ['JOB_NAME','GLUE_DB','SOURCE_TABLE','TEMP_BUCKET','DEST_BUCKET'])

job = Job(glueContext)
job.init(args['JOB_NAME'], args)

#pushing predicate aimed ate reducing cost by only fetchinng 48 hous of data
# the job will run by defautl every day and we include 1 extra da of data for redundancy 
# in case the job fails
today = datetime.now()
yesterday = today - timedelta(days=1)
todayYear=today.strftime('%Y')
todayMonth=today.strftime('%m')
todayDay=today.strftime('%d')
yesterdayYear=yesterday.strftime('%Y')
yesterdayMonth=yesterday.strftime('%m')
yesterdayDay=yesterday.strftime('%d')
#By default the partitions created by the crawler are named partition_0, partition_1, partition_2
#these are the names we need to use in our predicate push down expression
# partition_0 => year
# partition_1 => month
# partition_2 => day
#TODO: the bellow logic should be created by a function that take sthe # days in the past as param
today = "partition_0='"+todayYear+"' and partition_1='"+todayMonth+"' and partition_2='"+todayDay+"'"
yesterday = "partition_0=='"+yesterdayYear+"' and partition_1=='"+yesterdayMonth+"' and partition_2=='"+yesterdayDay+"'"
pdp= "(("+today+") or ("+yesterday+"))"

sensorsData = glueContext.create_dynamic_frame.from_catalog(database=args["GLUE_DB"], table_name=args["SOURCE_TABLE"], push_down_predicate = pdp)

dfc = sensorsData.relationalize("sensor_data_flat", "s3://"+args["TEMP_BUCKET"]+"/temp-dir/")

sensorsDataFlat = dfc.select('sensor_data_flat')
sensorsDataFlat = sensorsDataFlat.rename_field("partition_0", "year")
sensorsDataFlat = sensorsDataFlat.rename_field("partition_1", "month")
sensorsDataFlat = sensorsDataFlat.rename_field("partition_2", "day")
sensorsDataFlat = sensorsDataFlat.rename_field("partition_3", "hour")

#temporary workaround: drop here duplicate fields created by lowercasing
sensorsDataFlat = sensorsDataFlat.drop_fields(['lastadvertisement.datalen'])

#By default, spark dataframe overwites all data (even partitions that do not have new data). 
#We use the partitionOverwriteMode=dunamic to only overwrite new partitions.
spark = glueContext.spark_session
spark.conf.set('spark.sql.sources.partitionOverwriteMode','dynamic')
sensorsDataFlat.toDF().write.partitionBy("year","month","day","hour").mode("overwrite").format("parquet").save("s3://"+args["DEST_BUCKET"]+"/all")

job.commit()



