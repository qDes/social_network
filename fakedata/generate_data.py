import bcrypt

from mimesis import Person
from mimesis.enums import Gender

if __name__ == "__main__":
    person = Person('ru')

    for i in range(0, 30):
        a = person.full_name(gender=Gender.MALE)
        id_ = i + 1
        username = f"user{id_}"
        # password = bcrypt.hashpw(b'123456', bcrypt.gensalt())
        password = "fake"
        name, second_name = a.split()
        print(id_, password, name, second_name)
