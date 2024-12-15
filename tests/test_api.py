import pytest
from fastapi.testclient import TestClient
from app.main import app

client = TestClient(app)

def test_get_all_nodes():
    response = client.get("/nodes")
    assert response.status_code == 200
    assert isinstance(response.json(), list)

def test_get_node_relationships():
    response = client.get("/nodes/test-node")
    assert response.status_code in [200, 404]

def test_create_node():
    headers = {"X-API-Key": "your-secret-key"}
    node_data = {
        "id": "test-node",
        "label": "Test Label"
    }
    response = client.post("/nodes", json=node_data, headers=headers)
    assert response.status_code == 200

def test_delete_node():
    headers = {"X-API-Key": "your-secret-key"}
    response = client.delete("/nodes/test-node", headers=headers)
    assert response.status_code == 200 