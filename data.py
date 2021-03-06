import pandas as pd
import csv,json
import datetime
from datetime import datetime


costomer_company_df=pd.read_csv('test_data/Test task - Mongo - customer_companies.csv')
costomer_company_df['company_id'].isna()==False
costomer_df=pd.read_csv('test_data/Test task - Mongo - customers.csv')
orders_df=pd.read_csv('test_data/Test task - Postgres - orders.csv')
order_items_df=pd.read_csv('test_data/Test task - Postgres - order_items.csv')
deliveries_df=pd.read_csv('test_data/Test task - Postgres - deliveries.csv')

costomer_df=pd.merge(costomer_df,costomer_company_df,left_on='company_id',right_on='company_id',how='inner')
orders_df=pd.merge(orders_df,costomer_df,left_on='customer_id',right_on='user_id',how='inner')
order_items_df=pd.merge(order_items_df,orders_df,left_on='order_id',right_on='id',how='inner')
order_merged_df=pd.merge(order_items_df,deliveries_df,left_on='id_x',right_on='order_item_id',how='left')

delivered_amount = order_merged_df["delivered_quantity"] * order_merged_df["price_per_unit"]
order_merged_df["delivered_amount"] = round(delivered_amount,2)

total_amount = order_merged_df["quantity"] * order_merged_df["price_per_unit"]
order_merged_df["total_amount"] = round(total_amount, 2)




print(order_merged_df)
df = pd.DataFrame(order_merged_df, columns= ['order_id','created_at','order_name','id_x','price_per_unit','quantity','delivered_quantity','product','user_id','login','password','name','company_id','company_name','credit_cards','delivered_amount','total_amount'])
df.to_csv ('test_data.csv', index = False, header=False)
test_df = pd.read_csv('test_data.csv')


