import sys
import re
import random
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
# the job will run by defautl every day
#By default the partitions created by the crawler are named partition_0, partition_1, partition_2
#these are the names we need to use in our predicate push down expression
# partition_0 => year
# partition_1 => month
# partition_2 => day
today = datetime.now()
daysInPast = 2
predicates = []
for x in range(daysInPast):
    day = today - timedelta(days=x)
    dayYear=day.strftime('%Y')
    dayMonth=day.strftime('%m')
    dayDay=day.strftime('%d')
    predicates.append("partition_0='"+dayYear+"' and partition_1='"+dayMonth+"' and partition_2='"+dayDay+"'")

pdp= "(("+ ") OR (".join(predicates)+"))"
sensorsData = glueContext.create_dynamic_frame.from_catalog(database=args["GLUE_DB"], table_name=args["SOURCE_TABLE"], push_down_predicate = pdp)

dfc = sensorsData.relationalize("sensor_data_flat", "s3://"+args["TEMP_BUCKET"]+"/temp-dir/")

sensorsDataFlat = dfc.select('sensor_data_flat')
sensorsDataFlat = sensorsDataFlat.rename_field("partition_0", "year")
sensorsDataFlat = sensorsDataFlat.rename_field("partition_1", "month")
sensorsDataFlat = sensorsDataFlat.rename_field("partition_2", "day")
sensorsDataFlat = sensorsDataFlat.rename_field("partition_3", "hour")

#Renaming columns that will end up duplicates when lowercased
#thils function address the case where devices would have fieilds name the same 
#except for the caracter case in which case dupplicates would be created
columns = sensorsDataFlat.toDF().columns
#print(columns)
existing={}
toRename=[]
for col in columns:
    lowerCol = col.lower()
    #print(lowerCol)
    if not lowerCol in existing:
        existing[lowerCol] = 1
    else:
        toRename.append(col)
#print(toRename)
for col in toRename:
    newCol = col + str(random.randint(0,1000))
    while newCol.lower() in existing:
        newCol = col + str(random.randint(0,1000))
    #print("renaming" + col + " in " + newCol)
    sensorsDataFlat = sensorsDataFlat.rename_field(col, newCol)

#By default, spark dataframe overwites all data (even partitions that do not have new data). 
#We use the partitionOverwriteMode=dunamic to only overwrite new partitions.
spark = glueContext.spark_session
spark.conf.set('spark.sql.sources.partitionOverwriteMode','dynamic')
sensorsDataFlat.toDF().write.partitionBy("year","month","day","hour").mode("overwrite").format("parquet").save("s3://"+args["DEST_BUCKET"]+"/all")

job.commit()



