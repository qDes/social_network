c = box.schema.space.create('users')
c:format({
         {name = 'id', type = 'unsigned'},
         {name = 'username', type = 'string'},
         {name = 'first_name', type = 'string'},
         {name = 'second_name', type = 'string'}
         })
c:create_index('primary', {
         type = 'TREE',
         parts = {'id'}
         })
c:create_index('secondary', {
         type = 'TREE', unique=false,
         parts = {'first_name', 'second_name'}
         })