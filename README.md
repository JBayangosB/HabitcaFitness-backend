# HabticaFitness-backend
HabticaFitness was an attempt to use the RPG elements inspired by Habtica to create an application focussed on running.

#Rest API

website.appspot.com/Region/<Melbourne>
Returns quest particular to region

website.appspot.com/CreateUser/<Name>
Create user on server

website.appspot.com/ModifyUser/?name="Bob"&distance="5"
Adds distance ran to user

website.appspot.com/ModifyUser/?name="Bob"
Returns user

website.appspot.com/ModifyUser/?region=Melbourne&name=Bill&km=5
Attack monster in region and receive reward

#Objects

Quest
Name string
Desc string
HealthCurrent int64
HealthMax int64
MonsterName string
QuestLogs []QuestLog


QuestLog
Name string
Content string

User
Name string
HealthCurrent int64
HealthMax int64
ExpCurrent int64
ExpMax int64
Gold int64
DistanceTotal int64 
Level int64

Reward
Exp int64
Gold int64
