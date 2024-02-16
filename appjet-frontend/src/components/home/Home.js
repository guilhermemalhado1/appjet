// Import necessary modules
import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { AuthData } from "../../auth/AuthWrapper";
import ReactJson from 'react-json-view';
import './Home.css'; // Import the CSS file

// Define the Home component
function Home() {
  // Initialize state variables and functions
  const navigate = useNavigate();
  const { isAuthenticated, user, logout } = AuthData();
  const [showDocuments, setShowDocuments] = useState(false);
  const [showMonitoring, setShowMonitoring] = useState(false);
  const [appContainerState, setAppContainerState] = useState(false);
  const [databaseContainerState, setDatabaseContainerState] = useState(false);
  const [clusterData, setClusterData] = useState(null);
  const [clusterNames, setClusterNames] = useState([]);
  const [viewType, setViewType] = useState('json'); // Track view type
  const [folderStates, setFolderStates] = useState({
    folder1: false,
    subfolder1: false,
    subfolder2: false,
    subfolder3: false
  });

  // Effect hook for authentication
  useEffect(() => {
    if (!isAuthenticated) {
      navigate('/login');
    }
  }, [isAuthenticated, navigate]);

  // Effect hook for fetching cluster data
  useEffect(() => {
    if (showMonitoring) {
      fetchClusterData();
    }
  }, [showMonitoring]);

  // Function to fetch cluster data
  const fetchClusterData = async () => {
    try {
      const baseURL = process.env.REACT_APP_BASE_URL || "http://localhost:9999";
      const response = await fetch(baseURL + '/appjet/inspect', {
        headers: {
          Authorization: `${user.token}`
        }
      });
      const data = await response.json();
      setClusterData(data);
      if (Array.isArray(data['daemon-responses'])) {
        const arr = data['daemon-responses'].flatMap(cluster => Object.keys(cluster));
        setClusterNames(arr);
      }
    } catch (error) {
      console.error('Error fetching cluster data:', error);
    }
  };

  // Function to handle logout
  const handleLogout = async () => {
    await logout();
  };

  // Function to handle documentation button click
  const handleDocumentation = () => {
    setShowDocuments(!showDocuments);
    setShowMonitoring(false);
  };

  // Function to handle monitoring button click
  const handleMonitoring = () => {
    setShowMonitoring(!showMonitoring);
    setShowDocuments(false);
  };

  // Function to handle folder toggle
  const handleFolderToggle = (folder) => {
    setFolderStates(prevState => ({
      ...prevState,
      [folder]: !prevState[folder]
    }));
  };




// Function to handle view type change
const handleViewTypeChange = (type) => {
  let clusters = clusterData["daemon-responses"];
  clusters.forEach((cluster) => {
      Object.keys(cluster).forEach((elem1) => {
          Object.keys(cluster[elem1]).forEach((elem2) => {
            Object.keys(cluster[elem1][elem2]).forEach((elem3) => {
              let server = cluster[elem1][elem2][elem3]
              let appContainer = server.docker.app
              let dbContainer = server.docker.database
              
              alert("The container with name app, is " + (appContainer ? 'online' : 'offline') + ". The container with name database is " + (dbContainer ? 'online' : 'offline') + ".")
            });
          });
      });
  });
  setViewType(type); // Update the view type
};


  // Function to render nested rows
  const renderRows = (data) => {
    return Object.entries(data).map(([key, value]) => (
      <tr key={key}>
        <td>{key}</td>
        <td>{renderValue(value)}</td>
      </tr>
    ));
  };

  // Function to render nested values
  const renderValue = (value) => {
    if (typeof value === 'object' && value !== null) {
      if (Array.isArray(value)) {
        return (
          <ul>
            {value.map((item, index) => (
              <li key={index}>{renderValue(item)}</li>
            ))}
          </ul>
        );
      } else {
        return (
          <table className="nested-table">
            <tbody>
              {renderRows(value)}
            </tbody>
          </table>
        );
      }
    } else {
      return value;
    }
  };

  // JSX rendering
  return (
    <div className="home-container">
      <nav className="navbar">
        <div className="navbar-logo">AppJet | <span className="logged-in-as">Logged in as: {user.name}</span></div>
        <button className="logout-button" onClick={handleLogout}>Logout</button>
      </nav>
      <div className="buttons-container">
        <button className={`documentation-button ${showDocuments ? 'active' : ''}`} onClick={handleDocumentation}>
          <img src="https://cdn.pixabay.com/photo/2013/07/13/11/36/documents-158461_1280.png" alt="Documentation" />
          <span>Documentation</span>
        </button>
        <button className={`monitoring-button ${showMonitoring ? 'active' : ''}`} onClick={handleMonitoring}>
          <img src="https://cdn-icons-png.flaticon.com/512/3703/3703299.png" alt="Monitoring" />
          <span>Monitoring</span>
        </button>
      </div>
      {showDocuments && (
        <div className="additional-content">
           {showDocuments && (
        <div className="additional-content">
          {/* Tree style of windows file explorer with file icons */}
          <div className="tree">
            <ul>
            <li>
                <span className="folder" onClick={() => handleFolderToggle('folder1')}>
                  {folderStates.folder1 ? '▼' : '▶'} Appjet Documentation
                </span>
                {folderStates.folder1 && (
                  <ul>
                    <li>
                      <span className="folder" onClick={() => handleFolderToggle('subfolder1')}>
                        {folderStates.subfolder1 ? '▼' : '▶'} Subfolder 1
                      </span>
                      {folderStates.subfolder1 && (
                        <ul>
                          <li><a href='/file1.a'><span className="file">File 1.a</span></a></li>
                          <li><a href='/file1.b'><span className="file">File 1.b</span></a></li>
                        </ul>
                      )}
                    </li>
                    <li>
                      <span className="folder" onClick={() => handleFolderToggle('subfolder2')}>
                        {folderStates.subfolder2 ? '▼' : '▶'} Subfolder 2
                      </span>
                      {folderStates.subfolder2 && (
                        <ul>
                          <li><a href='/file2.a'><span className="file">File 2.a</span></a></li>
                          <li><a href='/file2.b'><span className="file">File 2.b</span></a></li>
                        </ul>
                      )}
                    </li>
                    <li>
                      <span className="folder" onClick={() => handleFolderToggle('subfolder3')}>
                        {folderStates.subfolder3 ? '▼' : '▶'} Subfolder 3
                      </span>
                      {folderStates.subfolder3 && (
                        <ul>
                          <li><a href='/file3.a'><span className="file">File 3.a</span></a></li>
                          <li><a href='/file3.b'><span className="file">File 3.b</span></a></li>
                        </ul>
                      )}
                    </li>
                  </ul>
                )}
              </li>
              
            </ul>
          </div>
        </div>
      )}
        </div>
      )}
      {showMonitoring && (
        <div className="additional-content">
          <div className="view-type-buttons">
            <button className={`view ${viewType === 'json' ? 'active' : ''}`} onClick={() => handleViewTypeChange('json')}>Check State</button>
          </div>
          {viewType === 'json' && clusterData && (
            <div>
              <h3>Appjet Configuration Details:</h3>
              <ReactJson src={clusterData} theme="rjv-default" />
            </div>
          )}
        </div>
      )}
    </div>
  );
}

export default Home;
