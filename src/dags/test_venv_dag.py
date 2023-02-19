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


def x():
	pass


with DAG(
	dag_id="example_python_venv_operator",
	schedule=None,
	start_date=pendulum.datetime(2021, 1, 1, tz="UTC"),
	catchup=False,
	tags=["example"],
) as dag:

	@task.virtualenv(
		task_id="virtualenv_python", requirements=["colorama==0.4.0"], system_site_packages=False
        )
	def callable_virtualenv():
		"""
		Example function that will be performed in a virtual environment.

		Importing at the module level ensures that it will not attempt to import the
		library before it is installed.
		"""
		from time import sleep
		from colorama import Back, Fore, Style
		print(Fore.RED + "some red text")
		print(Back.GREEN + "and with a green background")
		print(Style.DIM + "and in dim text")
		print(Style.RESET_ALL)
		for _ in range(4):
			print(Style.DIM + "Please wait...", flush=True)
			sleep(1)
		print("Finished")

	virtualenv_task = callable_virtualenv()
        
	virtualenv_task

	# [START howto_operator_python_venv_classic]
	virtual_classic = PythonVirtualenvOperator(
		task_id="virtualenv_classic",
	requirements="colorama==0.4.0",
	python_callable=x,
	)
	# [END howto_operator_python_venv_classic]

	virtual_classic
