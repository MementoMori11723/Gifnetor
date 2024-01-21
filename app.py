from flask import Flask,render_template
import ffmpeg

def convert(filename,output):
    stream = ffmpeg.input(filename)
    stream = ffmpeg.output(stream,output)
    ffmpeg.run_async(stream)

app = Flask(__name__)
@app.route("/")
def index():
    return render_template("index.html")
@app.route("/about/")
def about():
    return render_template("about.html")
