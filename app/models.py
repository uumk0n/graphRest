from pydantic import BaseModel
from typing import Optional, List, Dict

class VkUser(BaseModel):
    id: int
    screen_name: Optional[str] = None
    first_name: Optional[str] = None
    sex: Optional[str] = None
    city: Optional[Dict[str, str]] = None

class VkGroup(BaseModel):
    id: int
    name: str
    screen_name: str

class Relationship(BaseModel):
    source_id: int
    target_id: int
    relationship_type: str
    properties: Optional[Dict[str, any]] = None

class NodeWithRelationships(BaseModel):
    user: VkUser
    relationships: List[Relationship] = []