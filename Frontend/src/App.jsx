import { useState } from 'react'

import './App.css'
import {BrowserRouter, Navigate, Route, Routes} from 'react-router-dom'
import Index from './Pages/index'
import 'bootstrap/dist/css/bootstrap.min.css';


function App() {
  const [count, setCount] = useState(0)

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Index />} />
        {/*<Route path="/reportes" element={<Reportes />} />*/}
        <Route path="*" element={<Navigate to="/" replace={true} />} exact={true} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
