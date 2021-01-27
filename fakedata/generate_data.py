import bcrypt
import csv

from mimesis import Person
from mimesis.enums import Gender


if __name__ == "__main__":

    person = Person('en')
    password = bcrypt.hashpw(b'123456', bcrypt.gensalt()).decode()
    with open("fake_users.csv", mode='w') as my_file:
        user_writer = csv.writer(my_file, delimiter=',', quotechar='"', quoting=csv.QUOTE_MINIMAL)

        for i in range(0, 30):
            a = person.full_name(gender=Gender.MALE)
            id_ = i + 1
            username = f"user{id_}"

            name, second_name = a.split()
            print(username, password, name, second_name)
            user_writer.writerow([username, password, name, second_name])

