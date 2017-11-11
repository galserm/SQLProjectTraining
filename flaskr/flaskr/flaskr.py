import os
import sqlite3
from flask import Flask, request, session, g, redirect, url_for, abort, render_template, flash

import psycopg2

app = Flask(__name__)
app.config.from_object(__name__)


def connect_postgres_db():
    try:
        conn = psycopg2.connect("dbname='projectdb' user='neotek' host='localhost' password='kringstone' port='5432'")
        return conn
    except:
        print ("Unable to connect to database")

app.config.update(dict(
    DATABASE=os.path.join(app.root_path, 'flaskr_db'),
    SECRET_KEY='development key',
    USERNAME='admin',
    PASSWORD='default'
))
app.config.from_envvar('FLASKR_SETTINGS', silent=True)

def connect_db():
    rv = sqlite3.connect(app.config['DATABASE'])
    rv.row_factory = sqlite3.Row
    return rv

def get_db():
    if not hasattr(g, 'sqlite_db'):
        g.sqlite_db = connect_db()
    return g.sqlite_db

def get_postgres_db():
    if not hasattr(g, 'pg_db'):
        g.pg_db = connect_postgres_db()
        g.pg_db = g.pg_db.cursor()
    return g.pg_db

@app.teardown_appcontext
def close_db(error):
    if hasattr(g, 'sqlite_db'):
       g.sqlite_db.close()

def init_db():
    db = get_db()
    with app.open_resource('schema.sql', mode='r') as f:
        db.cursor().executescript(f.read())
    db.commit()

@app.cli.command('initdb')
def initdb_command():
    init_db()
    print('Initialized the database.')

@app.route('/')
def home():
    return redirect(url_for('login'))

@app.route('/acceuil')
def acceuil():
    return render_template('home.html', error='you are logged in')

@app.route('/account', methods=['POST', 'GET'])
def account():
    error = None
    if request.method == 'POST':
        db = get_postgres_db()
        if request.form['firstname']:
            db.execute("""UPDATE users SET firstname='%s' WHERE username='%s'""" % (request.form['firstname'], session['username']))
        if request.form['lastname']:
            db.execute("""UPDATE users SET lastname='%s' WHERE username='%s'""" % (request.form['lastname'], session['username']))
        if request.form['mail_address']:
            db.execute("""UPDATE users SET mail_address='%s' WHERE username='%s'""" % (request.form['mail_address'], session['username']))
        if request.form['username']:
            db.execute("""UPDATE users SET username='%s' WHERE username='%s'""" % (request.form['username'], session['username']))
            session['username'] = request.form['username']
    elif request.method == 'GET':
        db.execute("""SELECT * FROM users WHERE username='%s'""" % (session['username'])) 
            
            
@app.route('/add', methods=['POST'])
def add_entry():
    if not session.get('logged_in'):
        abort(401)
    db = get_db()
    db.execute('insert into entries (title, text) values (?, ?)',
                 [request.form['title'], request.form['text']])
    db.commit()
    flash('New entry was successfully posted')
    return redirect(url_for('show_entries'))

@app.route('/login', methods=['GET', 'POST'])
def login():
    error = None
    if request.method == 'POST':
        db = get_postgres_db()
        resp = db.execute("""select username, password from users where username='%s' and password='%s'""" % (request.form['username'], request.form['password']))
        entries = resp.fetchall()
        if not entries:
            error = 'Invalid username or password'
        else:
            session['logged_in'] = True
            session['username'] = request.form['username']
            flash('You were logged in')
            return redirect(url_for('acceuil'))
    return render_template('login.html', error=error)

@app.route('/logout')
def logout():
    session.pop('logged_in', None)
    flash('You were logged out')
    return redirect(url_for('login'))

@app.route('/register', methods=['GET', 'POST'])
def register():
    error = None
    if request.method == 'GET':
        return render_template('register.html')
    if not request.form['firstname']:
        error = 'Field marked with a * are mandatory'
    elif not request.form['lastname']:
        error = 'Field marked with a * are mandatory'
    elif not request.form['username']:
        error = 'Field marked with a * are mandatory'
    elif not request.form['password']:
        error = 'Field marked with a * are mandatory'
    if error != None:
        return render_template('register.html', error=error)
    db = get_postgres_db()
    db.execute("""INSERT INTO users (firstname, lastname, username, password) VALUES ('%s', '%s', '%s', '%s');""" % (request.form['firstname'], request.form['lastname'], request.form['username'], request.form['password']))
    flash('Account has been created')
    return redirect(url_for('login'))    

def get_user_id():
    error=none
    db = get_postgres_db()
    resp = db.execute("""SELECT id FROM users WHERE username='%s';""" % session['username'])
    rows = resp.fetchall()
    if rows:
        return rows['id']
    error = "user not found"
    return error

def getcomments():
    error = none
    db = get_postgres_db()
    resp = db.execute("""SELECT content,user_id,post,post_date,picture_path FROM posts;)""")
    row = 
@app.route('/addposts', method=['POST'])
def addposts():
    error=none
    db = get_postgres_db()
    resp = db.execute("""INSERT INTO posts (content, user_id, likes_number) VALUES ('%s', '%s', '0');""" % (request.form['content'], request.form['user']))
    row = resp.fetchall()
    if !row:
        error = "Failed to add new comment"
    return error
