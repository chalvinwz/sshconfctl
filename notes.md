Host d-gpay-fds
  HostName 10.11.15.210
  User chalvin.wiradhika
  PreferredAuthentications publickey
  IdentityFile ~/.ssh/chalvin

Host IDDR_Plane
  HostName 10.21.15.11
  User chalvin.wiradhika
  PreferredAuthentications publickey
  IdentityFile ~/.ssh/chalvin

Host Ansible
  HostName 10.11.15.15
  User chalvin.wiradhika
  PreferredAuthentications publickey
  IdentityFile ~/.ssh/chalvin

hey, i have an idea about a project using golang

i currently use lazyssh (https://github.com/Adembc/lazyssh), it reads & display servers from my  ~/.ssh/config and i need a way to add a new ssh config to that file. its such a repetitive task to do if i wanna add new entry of a server. i wanna create a cli app using go that does this:

user will asked for input: 
host
    - user input
hostname
    - user input
user 
    - set default user: chalvin.wiradhika
    - default user can be changed later
    - user can input and overwrite the default
    - user can just hit enter on empty and fallback to default user
prefferedaauthentication 
    - (no input, default to publickey)
identityfile (set default to ~/.ssh/chalvin, but user can change, and user can )
    - set default key path: ~/.ssh/chalvin
    - default key path can be changed later
    - user can input and overwrite the default
    - user can just hit enter on empty and fallback to default location

do best practice for this project, project structure, best stack, code spliting, logging, dockerfile

ask me what you need first or if you have something that is unclear