import os
import pytest
from testcontainers.postgres import PostgresContainer

from postgres import customers

pg = PostgresContainer("postgres:16-alpine")


@pytest.fixture(scope="module", autouse=True)
def setup(request):
    pg.start()

    def remove_container():
        pg.stop()

    # This will run after the test even if it throws an exception
    # guaranteeing that the data will be cleaned up after each test
    request.addfinalizer(remove_container)
    os.environ["DB_CONN"] = pg.get_connection_url()
    os.environ["DB_HOST"] = pg.get_container_host_ip()
    os.environ["DB_PORT"] = pg.get_exposed_port(5432)
    os.environ["DB_USERNAME"] = pg.username
    os.environ["DB_PASSWORD"] = pg.password
    os.environ["DB_NAME"] = pg.dbname
    customers.create_table()


@pytest.fixture(scope="function", autouse=True)
def setup_data():
    customers.delete_all_customers()


def test_create_customer():
    customers.create_customer("John", "jj@gmail.com")
    customer = customers.get_customer_by_email("jj@gmail.com")

    assert customer.name == "John"
    assert customer.email == "jj@gmail.com"


def test_get_all_customers():
    customers.create_customer("Siva", "siva@gmail.com")
    customers.create_customer("James", "james@gmail.com")
    customers_list = customers.get_all_customers()
    assert len(customers_list) == 2


def test_get_customer_by_email():
    customers.create_customer("John", "john@gmail.com")
    customer = customers.get_customer_by_email("john@gmail.com")
    assert customer.name == "John"
    assert customer.email == "john@gmail.com"
