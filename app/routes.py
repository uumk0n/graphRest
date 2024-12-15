from fastapi import APIRouter, HTTPException, Depends
from .database import Neo4jRepository
from .models import VkUser, Relationship, NodeWithRelationships
from typing import List
from config.config import settings

router = APIRouter()

def get_db():
    db = Neo4jRepository(
        settings.neo4j_uri,
        settings.neo4j_user,
        settings.neo4j_password
    )
    try:
        yield db
    finally:
        db.close()

@router.get("/nodes")
async def get_all_nodes(db: Neo4jRepository = Depends(get_db)):
    try:
        return db.get_all_nodes()
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.get("/nodes/{node_id}")
async def get_node_with_relationships(
    node_id: int, 
    db: Neo4jRepository = Depends(get_db)
):
    try:
        return db.get_node_relationships(node_id)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/nodes")
async def create_node(
    data: NodeWithRelationships, 
    db: Neo4jRepository = Depends(get_db)
):
    try:
        db.save_user_node(data.user)
        for rel in data.relationships:
            db.save_relationship(rel)
        return {"message": "Node and relationships created successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.delete("/nodes/{node_id}")
async def delete_node(
    node_id: int, 
    db: Neo4jRepository = Depends(get_db)
):
    try:
        db.delete_node(node_id)
        return {"message": "Node deleted successfully"}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))