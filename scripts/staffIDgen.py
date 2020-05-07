import uuid
import mysql.connector
from dotenv import load_dotenv
import os
import bcrypt

load_dotenv()

mydb = mysql.connector.connect(
  host=os.getenv("DB_HOST"),
  port=os.getenv("DB_PORT"),
  user=os.getenv("DB_USER"),
  passwd=os.getenv("DB_PASS"),
  db=os.getenv("DB_NAME"),
)

mycursor = mydb.cursor()
mycursor.execute("SELECT staff_id FROM staff")
myresult = mycursor.fetchall()

uid = str("staff_" + uuid.uuid4().hex)[:27]

for x in myresult:
    # print(x)
    if uid in x:
        print("match", x)
    else:
        print("not match", x)

print("This uid hasn't been used yet", uid)

firstName = input("Please first name: ")
LastName = input("Please last name: ")
Email = input("Email:")
Password = input("Password:")

print(uid,firstName,LastName,Email,Password)


password = Password.encode()
hashed = bcrypt.hashpw(password, bcrypt.gensalt(10))

print(hashed)

mycursor = mydb.cursor()

sql = "INSERT INTO staff (staff_id, first_name, last_name, email, password) VALUES (%s, %s, %s, %s, %s)"
val = (uid,firstName,LastName,Email,hashed)
mycursor.execute(sql, val)

mydb.commit()

print(mycursor.rowcount, "record inserted.")

