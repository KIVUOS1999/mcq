import requests
import pymongo
from bs4 import BeautifulSoup
import pymongo

table = "mcq_database"

def get_page(url, topic, collection, page):
    url = url +"?page="
    url_to_fetch = url+str(page)

    r = requests.get(url_to_fetch)
    content = r.content

    soup = BeautifulSoup(content, features="lxml")
    tables = soup.find_all("table")
    try:
        for table in tables:
            counter = 0
            dic = {}
            for row in table.find_all("tr"):
                for cell in row.find_all("td"):
                    text = (cell.text).strip()
                    if counter == 1:
                        question = text
                    elif counter == 3:
                        op1 = text
                    elif counter == 5:
                        op2 = text
                    elif counter == 7:
                        op3 = text
                    elif counter == 9:
                        op4 = text
                    elif counter == 11:
                        answer_raw = text.split("\n")[0].strip()
                        answer = answer_raw.split(": ")[1]
                    counter += 1

            dic["question"] = question
            dic["options"] = {"1": op1, "2":op2, "3":op3, "4":op4}
            dic["topic"] = topic
            for i, j in dic["options"].items():
                if j == answer:
                    dic["answer"] = i
            
            add_entry(collection, dic)
    except:
        pass

def create_client():
    client = pymongo.MongoClient("mongodb://localhost:27017")
    collection = client["mcq"]["questions"]

    return collection

def add_entry(collection, entry):
    collection = collection.insert_one(entry)

collection = create_client()

url = "https://mcqquestions.net/practice/dbms-mcqs"
total_page = 7
topic = "dbms"

for page in range(1,total_page+1):
    get_page(url, topic, collection, page)
