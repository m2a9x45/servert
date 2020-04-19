from dotenv import load_dotenv
from datetime import datetime
import os
import uuid
import mysql.connector
import time

load_dotenv()

# get orders expires dates

mydb = mysql.connector.connect(
  host=os.getenv("DB_HOST"),
  port=os.getenv("DB_PORT"),
  user=os.getenv("DB_USER"),
  passwd=os.getenv("DB_PASS"),
  db=os.getenv("DB_NAME"),
)

mycursor = mydb.cursor()
mycursor.execute("SELECT order_id, user_id, expires_at FROM orders")
myresult = mycursor.fetchall()

print("Current time : ", time.time())

curentUnixtime = time.time()

for x in myresult:
    # print(x)
    # print("expires at", x[2])
    ts = int(x[2])
    # print(datetime.utcfromtimestamp(ts).strftime('%d-%m-%Y'))
    expires = datetime.utcfromtimestamp(ts).strftime('%d-%m-%Y')
    # print(datetime.now().strftime('%d-%m-%Y'))
    current = datetime.now().strftime('%d-%m-%Y')
    if current == expires:
        # Will only occur if the order expires today
        print("Expires today : ", x)

    if curentUnixtime >= ts:
        # Will happen if the current time is greater or equal to the expiry time. Will triger for orders expired in the past
        print("Expired in the past : ", x)

    print("--------------------------------------------------------------")



# check if that date is today

# do something