#!/usr/bin/env python3
import sqlite3
import plotly.express as px
import os
import json

CFG_FILE = f"{os.environ['HOME']}/.config/cave-logger/config.json"
OUTDIR = "/tmp/cl"

def initDbConn():
    # Read config file first
    with open(CFG_FILE) as c:
        cfg = json.load(c)

    dbconn = sqlite3.connect(f"{os.environ['HOME']}/{cfg['database']['filename']}")
    dbconn.row_factory=sqlite3.Row

    return dbconn
# END

class App():
    def __init__(self):
        self.dbconn = initDbConn()

        # Check for outdir
        try:
          os.makedirs(OUTDIR, 0o755)
        except FileExistsError:
          print("Out directory exists")


    def analysePeople(self):
        print("Analysing People")
        

        cur = self.dbconn.cursor()
        cur.execute("""
        SELECT 
		people.id AS 'id',
		people.name AS 'name',
		(
                    SELECT COUNT(1)
                      FROM trip_groups
                     WHERE trip_groups.caverid = people.id
		) as 'count'
	FROM people
	ORDER BY name""")
        raw_data = cur.fetchall()

        # Flip the data
        data = {
            'names': [], 'names_count': [],
        }

        for row in raw_data:
            data['names'].append(row['name'])
            data['names_count'].append(row['count'])
            
        fig = px.bar(data, x='names', y='names_count')
        fig.show()
        
    def analyseClubs(self):
        print("Analysing Clubs")
        
        cur = self.dbconn.cursor()
        cur.execute("""
        SELECT 
		DISTINCT(people.club) AS 'club',
                (
                    SELECT COUNT(1) FROM people p WHERE p.club = people.club
                ) as 'count'
	FROM people
	ORDER BY club""")
        raw_data = cur.fetchall()

        # Flip the data
        data = {
            'clubs': [], 'clubs_count': [],
        }

        for row in raw_data:
            data['clubs'].append(row['club'])
            data['clubs_count'].append(row['count'])
            
        fig = px.bar(data, x='clubs', y='clubs_count')
        fig.show()


    def analyseCaves(self):
        print("Analysing Caves")

        cur = self.dbconn.cursor()
        cur.execute("""
        SELECT 
		DISTINCT(locations.name) AS 'cave',
                (
                    SELECT COUNT(1) FROM trips t WHERE t.caveid = locations.id
                ) as 'count'
	FROM locations
	ORDER BY cave""")
        raw_data = cur.fetchall()

        # Flip the data
        data = {
            'caves': [], 'caves_count': [],
        }

        for row in raw_data:
            data['caves'].append(row['cave'])
            data['caves_count'].append(row['count'])
            
        fig = px.bar(data, x='caves', y='caves_count')
        fig.show()

    def dbCloseConn():
        self.dbconn.close()
# END



if __name__ == '__main__':
    app = App()
    app.analyseClubs()
    app.analysePeople()
    app.analyseCaves()

    app.dbCloseConn()
