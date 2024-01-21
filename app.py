from flask import Flask,render_template,request,redirect
import ffmpeg
import os

app = Flask(__name__)

def fetch(filename):
    path = os.path.join('uploads/',filename)
    convert(path,os.path.join('uploads/','output.gif'))
    print("Complete")


def convert(filename,output):
    stream = ffmpeg.input(filename)
    stream = ffmpeg.output(stream,output)
    ffmpeg.run(stream)
    return

@app.route("/")
def index():
    return render_template("index.html")

@app.route("/about/")
def about():
    return render_template("about.html")

@app.route("/upload/",methods=['POST'])
def upload():
    file = request.files['video']
    file.save(os.path.join('uploads/',file.filename)) if file else None
    fetch(file.filename)
    return redirect("/")