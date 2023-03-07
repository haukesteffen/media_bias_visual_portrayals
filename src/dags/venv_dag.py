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


with DAG(
	dag_id="example_python_venv_operator",
	schedule=None,
	start_date=pendulum.datetime(2021, 1, 1, tz="UTC"),
	catchup=False,
	tags=["example"],
) as dag:

	@task.virtualenv(
		task_id="virtualenv_python", requirements=["face-recognition"], system_site_packages=True
		)
	def callable_virtualenv():
		import face_recognition
		print(f"face_recognition version {face_recognition.__version__}")

	virtualenv_task = callable_virtualenv()
        
	virtualenv_task
