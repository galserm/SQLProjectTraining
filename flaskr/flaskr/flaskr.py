
import os
import sqlite3
from flask import Flask, request, session, g, redirect, url_for, abort, render_template, flash

import psycopg2

app = Flask(__name__)
app.config.from_object(__name__)


def connect_postgres_db():
    try:
        conn = psycopg2.connect("dbname='projectdb' user='neotek' host='localhost' password='kringstone' port='5432'")
        conn.autocommit = True
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
    db = get_postgres_db()
    if request.method == 'POST':
        if request.form['firstname']:
            db.execute("""UPDATE users SET firstname='%s' WHERE username='%s'""" % (request.form['firstname'], session['username']))
        if request.form['lastname']:
            db.execute("""UPDATE users SET lastname='%s' WHERE username='%s'""" % (request.form['lastname'], session['username']))
        if request.form['email']:
            db.execute("""UPDATE users SET mail_address='%s' WHERE username='%s'""" % (request.form['email'], session['username']))
        if request.form['username']:
            db.execute("""UPDATE users SET username='%s' WHERE username='%s'""" % (request.form['username'], session['username']))
            session['username'] = request.form['username']
    try:
        db.execute("""SELECT * FrOM users WHERE username='%s'""" % (session['username']))
    except Exception as e:
        flash(e)
        return render_template('account.html')
    rows = db.fetchall()
    return render_template('account.html', rows=rows)
            
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
        try:
            db.execute("""select username, password from users where username='%s' and password='%s'""" % (request.form['username'], request.form['password']))
        except Exception as e:
            flash(e)
            return render_template("login.html")
        entries = db.fetchall()
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
    
    try:
        db.execute("""INSERT INTO users (firstname, lastname, mail_address, birth_date, username, password) VALUES ('%s', '%s', '%s', '%s', '%s', '%s');""" % (request.form['firstname'], request.form['lastname'], request.form['email'], request.form['birthdate'], request.form['username'], request.form['password']))
    except Exception as e:
        flash(e)
        return render_template('register.html')
    return redirect(url_for('login'))    

def get_user_id():
    error = None
    db = get_postgres_db()
    try:
        db.execute("""SELECT id FROM users WHERE username='%s';""" % session['username'])
    except Exception as e:
        flash(e)
        return "Couldn't get user for this id"
    rows = db.fetchall()
    if rows:
        return rows[0][0]
    error = "user not found"
    return error

@app.route('/followuser', methods=['POST'])
def follow_user():
    db = get_postgres_db()
    user_id = get_user_id()
    try:
        db.execute("""INSERT INTO followers (follower_user_id, followed_user_id) VALUES ('%d', '%d')""" % (user_id, request.form['user_id']))
    except Exception as e:
        flash(e)
        
@app.route('/like_post', methods=['POST'])
def like_post:
    db = get_postgres_db()
    user_id = get_user_id()
    try:
        db.execute("""INSERT INTO posts_likes (user_id, post_id) VALUES ('%d', '%d')""" % (user_id, request.form['post_id']))
        db.execute("""SELECT likes_number FROM posts WHERE id='%d'""" % request.form['post_id'])
        db.execute("""UPDATE posts SET likes_number="%d" WHERE id='%d'""" % request.form['post_id'])
        
    except Exception as e:
        flash(e)

@app.route('/post/<id>', methods=['GET'])
def post(id):
    error = None
    db = get_postgres_db()
    try:
        db.execute("""SELECT content,user_id,post_date,likes_number FROM posts WHERE id='%d'""" % id)
        postrow = db.fetchall()
        db.execute("""SELECT content,user_id,post_date,likes_number FROM comments WHERE post_id='%d'""" % id)
        commentsrow = db.fetchall()
    except Exception as e:
        flash(e)
    #return render_template('show_post_and_comments.html', postrow=postrow, commentsrow=commentsrow)

@app.route('/feed', methods=['GET'])
def feed():
    error = None
    posts = get_posts()
    if posts:
        return render_template('show_entries.html', posts=posts)
    error = "Failed to retrieve datas from server, try again later"
    return render_template('show_entries.html', error=error)

def get_username_from_id(id):
    db = get_postgres_db()
    try:
        db.execute("""SELECT username FROM users WHERE id='%d'""" % id)
        row = db.fetchall()
        return row[0][0]
    except Exception as e:
        flash(e)

def get_posts():
    db = get_postgres_db()
    try:
        db.execute("""SELECT content,user_id,post_date,picture_path,id FROM posts""")
    except Exception as e:
        flash(e)
        return None
    rows = db.fetchall()
    mylist = list(rows)
    i = 0;
    for item in mylist:
        new_item = list(item)
        new_item[1] = get_username_from_id(item[1])
        mylist[i] = new_item
        i = i + 1
    rows = tuple(mylist)
    return rows

@app.route('/user/<id>', methods=['GET'])
def user(id):
    db = get_postgres_db()
    try:
        db.execute("""SELECT username,birth_date,firstname,lastname FROM users WHERE id='%d'""" % id)
        userrows = db.fetchall()
        db.execute("""SELECT content,post_date,likes_number FROM posts WHERE user_id='%d'""" % id)
        postrows = db.fetchall()
    except Exception as e:
        flash(e)

@app.route('/addcomments', methods=['POST'])
def addcomment():
    db = get_postgres_db()
    user_id = get_user_id()
    try:
        db.execute("""INSERT INTO comments (content,user_id,post_id) VALUES (%s,%d,%d)""" % (request.form['content'], user_id, request.form['post_id']))
    except Exception as e:
        flash(e)

@app.route('/addposts', methods=['POST'])
def addposts():
    error = None
    db = get_postgres_db()
    user_id = get_user_id()
    try:
        db.execute("""INSERT INTO posts (content, user_id, likes_number) VALUES ('%s', '%s', '0');""" % (request.form['text'], user_id))
    except Exception as e:
        flash(e)
    return redirect(url_for('feed'))
