from airflow import DAG
from airflow.operators.bash import BashOperator
from datetime import datetime
default_args = {
    'owner': 'hauke',
    'depends_on_past': False,
    'start_date': datetime(2023, 2, 15),
    'retries': 0,
}
test_dag = DAG(
    'test_bash_script_dag',
    default_args=default_args,
    schedule_interval='30 * * * *'
)
# Define the BashOperator task
bash_task = BashOperator(
    task_id='bash_task_execute_script',
    bash_command='echo "Hello World!"',
    dag=test_dag
)
# Set task dependencies
bash_task
