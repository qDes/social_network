from faker import Faker
import random
import pymysql.cursors


if __name__ == "__main__":
    connection = pymysql.connect(host='127.0.0.1',
                                 port=6033,
                                 user='test',
                                 password='pzjqUkMnc7vfNHET',
                                 db='test',
                                 charset='utf8mb4',
                                 cursorclass=pymysql.cursors.DictCursor)

    fake = Faker()
    with connection.cursor() as cursor:
        for _ in range(100):
            id_user_1 = random.randint(1,100)
            id_user_2 = random.randint(1,100)
            id_lady_gaga = 200
            sql = "INSERT INTO messages (id, id_user_1, id_user_2, message) VALUES (uuid(), %s, %s, %s)"
            cursor.execute(sql, (id_user_1, id_user_2, fake.text()))
            connection.commit()
            #print(id_user_1, id_user_2, fake.text())
            #print(id_lady_gaga, id_user_2, fake.text())
'''
    with connection.cursor() as cursor:
        for i in range(n):
            country = fake.country()
            name = fake.name()
            city = fake.city()

            sql = "INSERT INTO customers (country, name, city) VALUES (%s, %s, %s)"
            cursor.execute(sql, (country, name, city))
            print(i)
            connection.commit()
'''