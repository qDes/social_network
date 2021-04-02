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
            sql = f"INSERT INTO messages (id, id_user_1, id_user_2, message) VALUES (uuid(), {id_user_1}, {id_user_2}, '{fake.text()}')"
            sql1 = f"INSERT INTO messages (id, id_user_1, id_user_2, message) VALUES (uuid(), 200, {id_user_2}, '{fake.text()}')"
            sql2 = f"INSERT INTO messages (id, id_user_1, id_user_2, message) VALUES (uuid(), {id_user_1}, 200, '{fake.text()}')"
            cursor.execute(sql1)
            connection.commit()
            cursor.execute(sql2)
            connection.commit()
            cursor.execute(sql)
            connection.commit()