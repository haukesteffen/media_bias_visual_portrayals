from __future__ import annotations

import logging
import shutil
import sys
import tempfile
import time
from pprint import pprint

import pendulum

from airflow import DAG
from airflow.decorators import task
from airflow.operators.python import ExternalPythonOperator, PythonVirtualenvOperator

log = logging.getLogger(__name__)

PATH_TO_PYTHON_BINARY = sys.executable

BASE_DIR = tempfile.gettempdir()

'''
def x():
	pass'''


with DAG(
	dag_id="example_python_venv_operator",
	schedule=None,
	start_date=pendulum.datetime(2021, 1, 1, tz="UTC"),
	catchup=False,
	tags=["example"],
) as dag:

	@task.virtualenv(
		task_id="virtualenv_python", requirements=["dlib", "face-recognition"], system_site_packages=False
        )
	def callable_virtualenv():
		import face_recognition
		print(f"face_recognition version {face_recognition.__version__}")

	virtualenv_task = callable_virtualenv()
        
	virtualenv_task

	'''# [START howto_operator_python_venv_classic]
	virtual_classic = PythonVirtualenvOperator(
		task_id="virtualenv_classic",
	requirements="colorama==0.4.0",
	python_callable=x,
	)
	# [END howto_operator_python_venv_classic]

	virtual_classic'''
