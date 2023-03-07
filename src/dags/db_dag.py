from datetime import datetime
from airflow import DAG
from airflow.decorators import task
from airflow.operators.python import PythonVirtualenvOperator

def get_face_encoding_from_image(img_path):
    # load image and create encoding
    image = face_recognition.load_image_file(img_path)
    encoding = face_recognition.face_encodings(image)
    return encoding

def get_face_encoding_from_bytes(bytearray):
    im = Image.open(io.BytesIO(bytearray))
    rgb = im.convert('RGB')
    array = np.array(rgb)
    encoding = face_recognition.face_encodings(array)
    return encoding

def check_for_faces(unknown_encoding, known_encodings):
    multiple_faces = len(unknown_encoding) > 1
    no_faces = len(unknown_encoding) == 0

    baerbock_encoding = known_encodings[0]
    contains_baerbock = (face_recognition.compare_faces([baerbock_encoding], unknown_encoding[0])[0])

    laschet_encoding = known_encodings[1]
    contains_laschet = (face_recognition.compare_faces([laschet_encoding], unknown_encoding[0])[0])

    scholz_encoding = known_encodings[2]
    contains_scholz = (face_recognition.compare_faces([scholz_encoding], unknown_encoding[0])[0])

    return multiple_faces, no_faces, contains_baerbock, contains_laschet, contains_scholz

with DAG(
	dag_id="example_db_dag",
	schedule=None,
	start_date=datetime(2023, 1, 1),
	catchup=False,
	tags=["example"],
) as dag:

	@task.virtualenv(
		task_id="virtualenv_python",
		requirements=["face-recognition", "sqlalchemy", "os", "pandas", "numpy", "PIL", "io"], 
		system_site_packages=True
		)
	def callable_virtualenv():
		# imports
		import face_recognition
		import pandas as pd
		import numpy as np
		from sqlalchemy import create_engine, text
		import os
		from PIL import Image
		import io


		# database connection
		username = os.environ["DBUSER"]
		password = os.environ["DBPASSWORD"]
		host = os.environ["DBHOSTNAME"]
		port = os.environ["DBPORT"]
		dbname = os.environ["DBDBNAME"]
		engine = create_engine(
			"postgresql://"
			+ username.strip()
			+ ":"
			+ password.strip()
			+ "@"
			+ host.strip()
			+ ":"
			+ port.strip()
			+"/"
			+ dbname.strip()
		)
		query = text("""SELECT * FROM items ORDER BY random() LIMIT 5""")

		# create dataframe with images
		with engine.begin() as con:
			df = pd.read_sql_query(sql=query, con=con)

		# get face encodings
		baerbock = get_face_encoding_from_image("data/training/annalena_baerbock.jpg")[0]
		laschet = get_face_encoding_from_image("data/training/armin_laschet.jpg")[0]
		scholz = get_face_encoding_from_image("data/training/olaf_scholz.jpg")[0]

		# move images to output directory if they contain a face
		df['multiple_faces'], df['no_faces'], df['baerbock'], df['laschet'], df['scholz'] = df.apply(lambda row: check_for_faces(get_face_encoding_from_bytes(row.data), [baerbock, laschet, scholz]), axis=1)
		print(df.head())
        

		# write back to db
		#todo

	virtualenv_task = callable_virtualenv()
        
	virtualenv_task
