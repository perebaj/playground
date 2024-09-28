import pytest
from sqlalchemy import Column, Integer, MetaData, String, Table, create_engine
from testcontainers.postgres import PostgresContainer


@pytest.fixture(scope="module")
def postgres_container():
    """Fixture to set up the PostgreSQL container."""
    with PostgresContainer("postgres:15") as postgres:
        yield postgres.get_connection_url()


def test_insert_into_table(postgres_container):
    """Test to insert data into a table and verify."""
    # Set up SQLAlchemy engine and metadata
    engine = create_engine(postgres_container)
    metadata = MetaData()

    # Define a new table
    test_table = Table("test_table", metadata, Column("id", Integer, primary_key=True), Column("name", String(50)))

    # Create the table in the database
    metadata.create_all(engine)

    # Insert an element into the table
    insert_statement = test_table.insert().values(id=1, name="Test Entry")

    with engine.connect() as connection:
        # Execute the insert statement
        connection.execute(insert_statement)

        # Query the table to check if the data was inserted correctly
        result = connection.execute(test_table.select()).fetchall()

    # Assert that the inserted data is correct
    assert len(result) == 1
    assert result[0] == (1, "Test Entry")
