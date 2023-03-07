#datetime
from datetime import timedelta, datetime

# The DAG object
from airflow import DAG

# Operators
from airflow.operators.python import PythonVirtualenvOperator

# initializing the default arguments
default_args = {
		'owner': 'hauke',
		'start_date': datetime(2023, 2, 15),
		'retries': 3,
		'retry_delay': timedelta(minutes=5)
}

# Instantiate a DAG object
hello_world_dag = DAG('hello_world_dag',
		default_args=default_args,
		description='Hello World DAG',
		schedule_interval='30 * * * *', 
		catchup=False,
		tags=['example, helloworld']
)

# python callable function
def print_hello():
		import numpy as np
		return 'Hello World!'

# Creating second task
hello_world_task = PythonVirtualenvOperator(
            task_id="test_task",
            requirements="numpy",
            python_callable=print_hello,
       	    dag=hello_world_dag)

# Set the order of execution of tasks. 
hello_world_task
