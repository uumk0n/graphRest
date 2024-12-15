from neo4j import GraphDatabase
from .models import VkUser, VkGroup, Relationship
from typing import List, Dict, Any

class Neo4jRepository:
    def __init__(self, uri: str, user: str, password: str):
        self.driver = GraphDatabase.driver(uri, auth=(user, password))

    def close(self):
        self.driver.close()

    def save_user_node(self, user: VkUser) -> None:
        with self.driver.session() as session:
            query = """
                MERGE (u:User {id: $id})
                ON CREATE SET u.screen_name = $screen_name, 
                            u.name = $name, 
                            u.sex = $sex, 
                            u.city = $city
                ON MATCH SET u.screen_name = $screen_name, 
                           u.name = $name, 
                           u.sex = $sex, 
                           u.city = $city
            """
            session.run(query, {
                "id": user.id,
                "screen_name": user.screen_name,
                "name": user.first_name,
                "sex": user.sex,
                "city": user.city.get("title") if user.city else None
            })

    def save_relationship(self, relationship: Relationship) -> None:
        with self.driver.session() as session:
            query = f"""
                MATCH (source:User {{id: $source_id}})
                MATCH (target:User {{id: $target_id}})
                MERGE (source)-[r:{relationship.relationship_type}]->(target)
            """
            session.run(query, {
                "source_id": relationship.source_id,
                "target_id": relationship.target_id
            })

    def get_all_nodes(self) -> List[Dict[str, Any]]:
        with self.driver.session() as session:
            result = session.run("MATCH (n:User) RETURN n")
            return [dict(record["n"].items()) for record in result]

    def get_node_relationships(self, node_id: int) -> List[Dict[str, Any]]:
        with self.driver.session() as session:
            query = """
                MATCH (n:User {id: $node_id})-[r]->(m)
                RETURN n as source, type(r) as type, m as target
            """
            result = session.run(query, {"node_id": node_id})
            return [{
                "source": dict(record["source"].items()),
                "type": record["type"],
                "target": dict(record["target"].items())
            } for record in result]

    def delete_node(self, node_id: int) -> None:
        with self.driver.session() as session:
            query = "MATCH (n:User {id: $node_id}) DETACH DELETE n"
            session.run(query, {"node_id": node_id})