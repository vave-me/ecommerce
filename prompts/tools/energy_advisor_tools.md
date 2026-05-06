<!-- Version: 1.0.0 -->
<!-- Last Updated: 2025-07-21 -->
<!-- Status: Active -->
<!-- Purpose: Energy advisor specific tool definitions -->

# Energy Advisor Tool Schemas

## 🔋 Core Energy Tools

### 1. Energy Consumption Analyzer
```json
{
  "name": "analyze_energy_consumption",
  "description": "Analyze current energy consumption patterns and costs",
  "parameters": {
    "type": "object",
    "properties": {
      "monthly_bill_pln": {
        "type": "number",
        "description": "Monthly electricity bill in PLN"
      },
      "annual_consumption_kwh": {
        "type": "number",
        "description": "Annual energy consumption in kWh"
      },
      "tariff_type": {
        "type": "string",
        "enum": ["G11", "G12", "G12w", "C11", "C12a", "C12b"],
        "description": "Energy tariff type"
      },
      "heating_type": {
        "type": "string",
        "enum": ["electric", "gas", "oil", "coal", "heat_pump", "other"],
        "description": "Current heating system type"
      },
      "building_area_m2": {
        "type": "number",
        "description": "Building area in square meters"
      }
    },
    "required": ["monthly_bill_pln"]
  }
}
```

### 2. PV System Calculator
```json
{
  "name": "calculate_pv_system",
  "description": "Calculate optimal photovoltaic system configuration",
  "parameters": {
    "type": "object",
    "properties": {
      "annual_consumption_kwh": {
        "type": "number",
        "description": "Annual electricity consumption in kWh"
      },
      "location": {
        "type": "object",
        "properties": {
          "city": {"type": "string"},
          "latitude": {"type": "number"},
          "longitude": {"type": "number"}
        }
      },
      "roof_type": {
        "type": "string",
        "enum": ["flat", "pitched", "east_west", "ground_mount"],
        "description": "Installation surface type"
      },
      "roof_angle_degrees": {
        "type": "number",
        "description": "Roof angle in degrees (0-90)"
      },
      "roof_orientation": {
        "type": "string",
        "enum": ["S", "SE", "SW", "E", "W", "NE", "NW", "N"],
        "description": "Roof orientation"
      },
      "available_area_m2": {
        "type": "number",
        "description": "Available installation area"
      },
      "shading_factor": {
        "type": "number",
        "minimum": 0,
        "maximum": 1,
        "description": "Shading factor (0=full shade, 1=no shade)"
      }
    },
    "required": ["annual_consumption_kwh", "location"]
  }
}
```

### 3. ROI Calculator
```json
{
  "name": "calculate_roi",
  "description": "Calculate return on investment for energy solutions",
  "parameters": {
    "type": "object",
    "properties": {
      "system_type": {
        "type": "string",
        "enum": ["pv", "heat_pump", "pv_heat_pump", "battery", "full_system"]
      },
      "investment_cost_pln": {
        "type": "number",
        "description": "Total investment cost in PLN"
      },
      "monthly_savings_pln": {
        "type": "number",
        "description": "Expected monthly savings in PLN"
      },
      "subsidies": {
        "type": "object",
        "properties": {
          "moj_prad": {"type": "number", "description": "Mój Prąd subsidy amount"},
          "czyste_powietrze": {"type": "number", "description": "Czyste Powietrze subsidy"},
          "ulga_termo": {"type": "number", "description": "Thermomodernization relief"},
          "other": {"type": "number", "description": "Other subsidies"}
        }
      },
      "energy_price_increase_percent": {
        "type": "number",
        "default": 5,
        "description": "Annual energy price increase percentage"
      },
      "system_degradation_percent": {
        "type": "number",
        "default": 0.5,
        "description": "Annual system efficiency degradation"
      }
    },
    "required": ["system_type", "investment_cost_pln", "monthly_savings_pln"]
  }
}
```

### 4. Heat Pump Selector
```json
{
  "name": "select_heat_pump",
  "description": "Select optimal heat pump based on building parameters",
  "parameters": {
    "type": "object",
    "properties": {
      "building_area_m2": {
        "type": "number",
        "description": "Building area to heat"
      },
      "insulation_quality": {
        "type": "string",
        "enum": ["poor", "average", "good", "excellent"],
        "description": "Building insulation quality"
      },
      "current_heating": {
        "type": "object",
        "properties": {
          "type": {"type": "string"},
          "annual_cost_pln": {"type": "number"},
          "fuel_consumption": {"type": "number"}
        }
      },
      "heat_pump_type": {
        "type": "string",
        "enum": ["air_to_water", "ground_source", "water_source"],
        "description": "Preferred heat pump type"
      },
      "hot_water_demand_liters": {
        "type": "number",
        "description": "Daily hot water demand"
      },
      "lowest_temp_celsius": {
        "type": "number",
        "description": "Lowest expected outside temperature"
      }
    },
    "required": ["building_area_m2", "insulation_quality"]
  }
}
```

### 5. Battery Storage Optimizer
```json
{
  "name": "optimize_battery_storage",
  "description": "Optimize battery storage capacity for PV system",
  "parameters": {
    "type": "object",
    "properties": {
      "pv_system_kwp": {
        "type": "number",
        "description": "PV system power in kWp"
      },
      "daily_consumption_profile": {
        "type": "array",
        "items": {
          "type": "number"
        },
        "description": "24-hour consumption profile in kWh"
      },
      "grid_sellback_rate": {
        "type": "number",
        "description": "Rate for selling energy back to grid (PLN/kWh)"
      },
      "time_of_use_tariff": {
        "type": "boolean",
        "description": "Whether time-of-use tariff is active"
      },
      "backup_power_hours": {
        "type": "number",
        "description": "Required backup power duration in hours"
      },
      "budget_pln": {
        "type": "number",
        "description": "Available budget for battery system"
      }
    },
    "required": ["pv_system_kwp", "daily_consumption_profile"]
  }
}
```

### 6. Installation Planner
```json
{
  "name": "plan_installation",
  "description": "Create detailed installation plan with timeline",
  "parameters": {
    "type": "object",
    "properties": {
      "system_components": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "type": {"type": "string"},
            "capacity": {"type": "number"},
            "quantity": {"type": "number"}
          }
        }
      },
      "installation_address": {
        "type": "object",
        "properties": {
          "street": {"type": "string"},
          "city": {"type": "string"},
          "postal_code": {"type": "string"},
          "province": {"type": "string"}
        }
      },
      "building_type": {
        "type": "string",
        "enum": ["single_family", "multi_family", "commercial", "industrial"]
      },
      "permits_required": {
        "type": "array",
        "items": {
          "type": "string",
          "enum": ["building_permit", "grid_connection", "environmental", "zoning"]
        }
      },
      "preferred_start_date": {
        "type": "string",
        "format": "date"
      }
    },
    "required": ["system_components", "installation_address", "building_type"]
  }
}
```

### 7. Subsidy Checker
```json
{
  "name": "check_available_subsidies",
  "description": "Check available government subsidies and grants",
  "parameters": {
    "type": "object",
    "properties": {
      "applicant_type": {
        "type": "string",
        "enum": ["individual", "business", "municipality", "housing_coop"]
      },
      "planned_investment": {
        "type": "array",
        "items": {
          "type": "string",
          "enum": ["pv_panels", "heat_pump", "insulation", "windows", "ventilation", "battery"]
        }
      },
      "building_age_years": {
        "type": "number",
        "description": "Age of the building"
      },
      "household_income_pln": {
        "type": "number",
        "description": "Annual household income (for income-based programs)"
      },
      "location_province": {
        "type": "string",
        "description": "Province for regional programs"
      }
    },
    "required": ["applicant_type", "planned_investment"]
  }
}
```

### 8. Energy Flow Simulator
```json
{
  "name": "simulate_energy_flow",
  "description": "Simulate energy production and consumption flow",
  "parameters": {
    "type": "object",
    "properties": {
      "simulation_period": {
        "type": "string",
        "enum": ["day", "week", "month", "year"],
        "description": "Period to simulate"
      },
      "pv_capacity_kwp": {
        "type": "number",
        "description": "PV system capacity"
      },
      "battery_capacity_kwh": {
        "type": "number",
        "description": "Battery storage capacity"
      },
      "consumption_data": {
        "type": "array",
        "items": {
          "type": "object",
          "properties": {
            "timestamp": {"type": "string"},
            "consumption_kw": {"type": "number"}
          }
        }
      },
      "weather_scenario": {
        "type": "string",
        "enum": ["sunny", "cloudy", "mixed", "historical"],
        "description": "Weather conditions for simulation"
      }
    },
    "required": ["simulation_period", "pv_capacity_kwp"]
  }
}
```

## 📊 Output Schemas

### Energy Analysis Result
```json
{
  "type": "object",
  "properties": {
    "current_state": {
      "annual_consumption_kwh": "number",
      "annual_cost_pln": "number",
      "co2_emissions_kg": "number"
    },
    "recommended_solution": {
      "components": "array",
      "total_cost_pln": "number",
      "annual_savings_pln": "number",
      "roi_years": "number",
      "co2_reduction_kg": "number"
    },
    "financial_analysis": {
      "investment_cost": "number",
      "available_subsidies": "number",
      "net_cost": "number",
      "payback_period_years": "number",
      "20_year_savings": "number"
    },
    "next_steps": "array"
  }
}
```

### Installation Plan Result
```json
{
  "type": "object",
  "properties": {
    "timeline": {
      "total_duration_days": "number",
      "phases": [
        {
          "name": "string",
          "duration_days": "number",
          "requirements": "array"
        }
      ]
    },
    "permits": {
      "required": "array",
      "estimated_time_days": "number",
      "estimated_cost_pln": "number"
    },
    "installation_details": {
      "mounting_type": "string",
      "cable_length_m": "number",
      "required_space_m2": "number"
    }
  }
}
```

## 🔧 Integration Notes

1. **Weather Data**: Integrate with local weather APIs for accurate solar predictions
2. **Tariff Updates**: Connect to energy provider APIs for current pricing
3. **Equipment Database**: Maintain updated database of panels, inverters, heat pumps
4. **Subsidy Portal**: Real-time connection to government subsidy programs
5. **Grid Requirements**: Local DSO (Distribution System Operator) requirements

## 🚀 Usage Examples

### Example 1: Complete Energy Audit
```javascript
// Step 1: Analyze current consumption
const consumption = await analyze_energy_consumption({
  monthly_bill_pln: 450,
  tariff_type: "G11",
  building_area_m2: 150
});

// Step 2: Calculate PV system
const pv_system = await calculate_pv_system({
  annual_consumption_kwh: consumption.annual_kwh,
  location: { city: "Warszawa", latitude: 52.23, longitude: 21.01 },
  roof_type: "pitched",
  roof_angle_degrees: 35
});

// Step 3: Check subsidies
const subsidies = await check_available_subsidies({
  applicant_type: "individual",
  planned_investment: ["pv_panels"]
});

// Step 4: Calculate ROI
const roi = await calculate_roi({
  system_type: "pv",
  investment_cost_pln: pv_system.total_cost,
  monthly_savings_pln: pv_system.estimated_savings,
  subsidies: subsidies
});
```

### Example 2: Heat Pump Replacement
```javascript
// Analyze heating needs and select optimal heat pump
const heat_pump = await select_heat_pump({
  building_area_m2: 200,
  insulation_quality: "good",
  current_heating: {
    type: "gas_boiler",
    annual_cost_pln: 6000
  },
  lowest_temp_celsius: -20
});
```