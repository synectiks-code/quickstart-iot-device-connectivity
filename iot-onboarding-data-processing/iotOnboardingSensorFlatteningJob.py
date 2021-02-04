import sys
import re
from awsglue.transforms import *
from awsglue.utils import getResolvedOptions
from pyspark.context import SparkContext
from awsglue.context import GlueContext
from awsglue.job import Job
from awsglue.dynamicframe import DynamicFrame

glueContext = GlueContext(SparkContext.getOrCreate())

args = getResolvedOptions(sys.argv, ['JOB_NAME','GLUE_DB','SOURCE_TABLE','TEMP_BUCKET','DEST_BUCKET'])

job = Job(glueContext)
job.init(args['JOB_NAME'], args)

sensorsData = glueContext.create_dynamic_frame.from_catalog(database=args["GLUE_DB"], table_name=args["SOURCE_TABLE"])

dfc = sensorsData.relationalize("sensor_data_flat", "s3://"+args["TEMP_BUCKET"]+"/temp-dir/")

sensorsDataFlat = dfc.select('sensor_data_flat')
sensorsDataFlat = sensorsDataFlat.rename_field("partition_0", "year")
sensorsDataFlat = sensorsDataFlat.rename_field("partition_1", "month")
sensorsDataFlat = sensorsDataFlat.rename_field("partition_2", "day")
sensorsDataFlat = sensorsDataFlat.rename_field("partition_3", "hour")

sensorsDataFlat.toDF().write.partitionBy("year","month","day","hour").mode("overwrite").format("parquet").save("s3://"+args["DEST_BUCKET"]+"/all")

job.commit()



