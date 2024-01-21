from app import app
from flask_frozen import Freezer
freezer = Freezer(app)
freezer.freeze()
# app.run(debug=True)