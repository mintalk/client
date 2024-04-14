# List of server commands
## User management

~~~
user add <name>
~~~

Create a new user. When he first logs in, a password will be created.

~~~
user op <user>
~~~

Make user an operator.

~~~
user deop <user>
~~~

Make user not an operator.
> **Note**: Operators are currently not implemented.

~~~
user del <name>
~~~

Remove a user.

~~~
user list
~~~

List all users.

## Group management

~~~
group add <name>
~~~

Create a new group.

~~~
group del <name>
~~~

Remove a group.

~~~
group move <group> <parent>
~~~

Make a group a child of another group.

~~~
group root <group>
~~~

Make group a root group (no parent).

~~~
group list
~~~

Show a tree of all groups.

## Channel management

~~~
channel add <name> <group>
~~~

Create a new channel as a child of a group.

~~~
channel del <name>
~~~

Remove a channel.

~~~
channel move <channel> <group>
~~~

Make a channel a child of a group.

~~~
channel list
~~~

List all channels.
