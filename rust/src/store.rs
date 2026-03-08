use std::collections::{HashMap, HashSet, VecDeque};

use crate::errors::{EventGraphError, Result};
use crate::event::Event;
use crate::types::{ActorId, ConversationId, EventId, EventType};

pub struct ChainVerification {
    pub valid: bool,
    pub length: usize,
}

pub trait Store {
    fn append(&mut self, event: Event) -> Result<Event>;
    fn get(&self, event_id: &EventId) -> Result<&Event>;
    fn head(&self) -> Option<&Event>;
    fn count(&self) -> usize;
    fn verify_chain(&self) -> ChainVerification;
    fn close(&mut self);
}

pub struct InMemoryStore {
    events: Vec<Event>,
    index: HashMap<String, usize>,
}

impl InMemoryStore {
    pub fn new() -> Self {
        Self { events: Vec::new(), index: HashMap::new() }
    }

    pub fn recent(&self, limit: usize) -> Vec<&Event> {
        self.events.iter().rev().take(limit).collect()
    }

    pub fn by_type(&self, event_type: &EventType, limit: usize) -> Vec<&Event> {
        self.events
            .iter()
            .rev()
            .filter(|e| e.event_type.value() == event_type.value())
            .take(limit)
            .collect()
    }

    pub fn by_source(&self, source: &ActorId, limit: usize) -> Vec<&Event> {
        self.events
            .iter()
            .rev()
            .filter(|e| e.source.value() == source.value())
            .take(limit)
            .collect()
    }

    pub fn by_conversation(&self, conversation_id: &ConversationId, limit: usize) -> Vec<&Event> {
        self.events
            .iter()
            .rev()
            .filter(|e| e.conversation_id.value() == conversation_id.value())
            .take(limit)
            .collect()
    }

    pub fn ancestors(&self, event_id: &EventId, max_depth: usize) -> Result<Vec<&Event>> {
        // Verify starting event exists
        if !self.index.contains_key(event_id.value()) {
            return Err(EventGraphError::EventNotFound {
                event_id: event_id.value().to_string(),
            });
        }

        let mut result = Vec::new();
        let mut visited = HashSet::new();
        let mut queue = VecDeque::new();

        // Seed with the causes of the starting event
        let start = &self.events[self.index[event_id.value()]];
        for cause_id in start.causes.iter() {
            if cause_id.value() != event_id.value() && visited.insert(cause_id.value().to_string()) {
                queue.push_back((cause_id.value().to_string(), 1usize));
            }
        }

        while let Some((id, depth)) = queue.pop_front() {
            if depth > max_depth {
                continue;
            }
            if let Some(&idx) = self.index.get(&id) {
                let event = &self.events[idx];
                result.push(event);

                if depth < max_depth {
                    for cause_id in event.causes.iter() {
                        if visited.insert(cause_id.value().to_string()) {
                            queue.push_back((cause_id.value().to_string(), depth + 1));
                        }
                    }
                }
            }
        }

        Ok(result)
    }

    pub fn descendants(&self, event_id: &EventId, max_depth: usize) -> Result<Vec<&Event>> {
        // Verify starting event exists
        if !self.index.contains_key(event_id.value()) {
            return Err(EventGraphError::EventNotFound {
                event_id: event_id.value().to_string(),
            });
        }

        let mut result = Vec::new();
        let mut visited = HashSet::new();
        let mut queue = VecDeque::new();
        queue.push_back((event_id.value().to_string(), 0usize));
        visited.insert(event_id.value().to_string());

        while let Some((id, depth)) = queue.pop_front() {
            if depth >= max_depth {
                continue;
            }

            // Find events that reference `id` in their causes
            for event in &self.events {
                let references_id = event.causes.iter().any(|c| c.value() == id);
                if references_id && visited.insert(event.id.value().to_string()) {
                    result.push(event);
                    queue.push_back((event.id.value().to_string(), depth + 1));
                }
            }
        }

        Ok(result)
    }
}

impl Default for InMemoryStore {
    fn default() -> Self { Self::new() }
}

impl Store for InMemoryStore {
    fn append(&mut self, event: Event) -> Result<Event> {
        if !self.events.is_empty() {
            let last = self.events.last().unwrap();
            if event.prev_hash.value() != last.hash.value() {
                return Err(EventGraphError::ChainIntegrity {
                    position: self.events.len(),
                    detail: format!(
                        "prev_hash {} != head hash {}",
                        event.prev_hash.value(),
                        last.hash.value()
                    ),
                });
            }
        }
        let id_str = event.id.value().to_string();
        self.events.push(event);
        let idx = self.events.len() - 1;
        self.index.insert(id_str, idx);
        Ok(self.events[idx].clone())
    }

    fn get(&self, event_id: &EventId) -> Result<&Event> {
        self.index
            .get(event_id.value())
            .map(|&i| &self.events[i])
            .ok_or_else(|| EventGraphError::EventNotFound {
                event_id: event_id.value().to_string(),
            })
    }

    fn head(&self) -> Option<&Event> {
        self.events.last()
    }

    fn count(&self) -> usize {
        self.events.len()
    }

    fn verify_chain(&self) -> ChainVerification {
        for i in 1..self.events.len() {
            if self.events[i - 1].hash.value() != self.events[i].prev_hash.value() {
                return ChainVerification { valid: false, length: i };
            }
        }
        ChainVerification { valid: true, length: self.events.len() }
    }

    fn close(&mut self) {}
}
