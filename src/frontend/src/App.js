import React, { useState, useEffect } from 'react';
import axios from 'axios';

function App() {
  const [zone, setZone] = useState('');
  const [records, setRecords] = useState([]);
  const [selected, setSelected] = useState(new Set());
  const [selectAll, setSelectAll] = useState(false);
  const [filter, setFilter] = useState('All');
  const [types, setTypes] = useState([]);

  const fetchRecords = async () => {
    const result = await axios.get(`http://localhost:8080/api/records/${zone}`);
    setRecords(result.data);
    setSelected(new Set());  // Clear selections when new records are fetched
    setSelectAll(false);  // Reset the selectAll state
    fetchTypes(); // Fetch types after fetching records
  };

  const fetchTypes = async () => {
    if (zone) {
      const result = await axios.get(`http://localhost:8080/api/types/${zone}`);
      setTypes(result.data);
    }
  };

  const handleDelete = async () => {
    for (let id of Array.from(selected)) {
      await axios.delete(`http://localhost:8080/api/record/${zone}/${id}`);
    }
    await fetchRecords();
  };

  const handleCheckboxChange = (event) => {
    const id = event.target.name;
    const newSelected = new Set(selected);
    if (selected.has(id)) {
      newSelected.delete(id);
    } else {
      newSelected.add(id);
    }
    setSelected(newSelected);
  };

  const handleSelectAllChange = () => {
    if (selectAll) {
      setSelected(new Set());
      setSelectAll(false);
    } else {
      setSelected(new Set(records.map(record => record.id)));
      setSelectAll(true);
    }
  };

  // Update the select all state when all records are manually selected
  useEffect(() => {
    if (records.length > 0 && records.every(record => selected.has(record.id))) {
      setSelectAll(true);
    } else {
      setSelectAll(false);
    }
  }, [records, selected]);

  const filteredRecords = records.filter(record => filter === 'All' || record.type === filter);

  return (
      <div>
        <input value={zone} onChange={(e) => setZone(e.target.value)} placeholder="Enter zone name" />
        <button onClick={fetchRecords}>Fetch Records</button>
        <button onClick={handleDelete}>Delete Selected</button>
        <button onClick={handleSelectAllChange}>{selectAll ? 'Deselect All' : 'Select All'}</button>
        <select value={filter} onChange={(e) => setFilter(e.target.value)}>
          <option value="All">All</option>
          {types.map((type, index) => (
              <option value={type} key={index}>{type}</option>
          ))}
        </select>
        <div>
          {filteredRecords.map((record) => (
              <div key={record.id}>
                <label>
                  <input
                      type="checkbox"
                      name={record.id}
                      checked={selected.has(record.id)}
                      onChange={handleCheckboxChange}
                  />
                  {record.name} ({record.type}): {record.content}
                </label>
              </div>
          ))}
        </div>
      </div>
  );
}

export default App;