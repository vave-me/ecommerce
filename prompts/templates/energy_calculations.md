<!-- Version: 1.0.0 -->
<!-- Last Updated: 2025-07-21 -->
<!-- Status: Active -->
<!-- Purpose: Energy calculation templates for AI Energy Advisor -->

# Energy Calculation Templates

## 📊 Solar PV Calculations

### Basic System Sizing
```
Formula: System Size (kW) = Annual Consumption (kWh) / (365 × Peak Sun Hours × System Efficiency)

Example:
- Annual consumption: 12,000 kWh
- Peak sun hours: 4.5 (average)
- System efficiency: 0.80
- Result: 12,000 / (365 × 4.5 × 0.80) = 9.1 kW system
```

### Monthly Production Estimate
```
Monthly Production (kWh) = System Size (kW) × Daily Sun Hours × Days × Performance Ratio

Performance Factors:
- Temperature coefficient: -0.4%/°C above 25°C
- Inverter efficiency: 97-98%
- Wiring losses: 2%
- Soiling losses: 2-5%
- Shading losses: Variable
- Total performance ratio: 0.75-0.85
```

### Financial Calculations
```
Simple Payback = Net System Cost / Annual Savings

Where:
- Net System Cost = Total Cost - Incentives - Tax Credits
- Annual Savings = (kWh Produced × Electricity Rate) + Avoided Costs

NPV = Σ(Cash Flow_t / (1 + r)^t) - Initial Investment
IRR = Rate where NPV = 0
LCOE = (Total Lifetime Cost) / (Total Lifetime Production)
```

## 🔥 Heat Pump Calculations

### Heating Load Calculation
```
Heat Loss (BTU/hr) = Building Area × (Indoor Temp - Outdoor Temp) × Heat Loss Factor

Heat Loss Factors (BTU/hr/sq ft/°F):
- Poor insulation: 35-50
- Average insulation: 25-35
- Good insulation: 15-25
- Excellent insulation: 10-15

Heat Pump Size (tons) = Heat Loss (BTU/hr) / 12,000
```

### COP and Efficiency
```
COP (Coefficient of Performance) = Heat Output / Energy Input

Seasonal Performance:
- HSPF (Heating): 8.5-13+ (higher is better)
- SEER (Cooling): 15-25+ (higher is better)

Annual Heating Cost = (Annual Heat Load × Cost per kWh) / (HSPF × 3.412)
```

### Cost Comparison
```
Current System Annual Cost:
- Gas: (Therms × $/Therm) / Furnace Efficiency
- Oil: (Gallons × $/Gallon) / Furnace Efficiency
- Electric: kWh × $/kWh

Heat Pump Annual Cost:
- Heating: (Annual BTU / HSPF) × $/kWh × 0.000293
- Cooling: (Annual BTU / SEER) × $/kWh × 0.000293

Annual Savings = Current Cost - Heat Pump Cost
```

## 🔋 Battery Storage Calculations

### Battery Sizing
```
Required Capacity (kWh) = Daily Usage (kWh) × Days of Autonomy × Depth of Discharge

Example:
- Daily usage: 30 kWh
- Autonomy desired: 1 day
- DoD (LiFePO4): 90%
- Required: 30 × 1 / 0.90 = 33.3 kWh

Round-trip Efficiency:
- Lithium-ion: 90-95%
- Lead-acid: 80-85%
```

### Economic Analysis
```
Daily Arbitrage Value = (Peak Rate - Off-Peak Rate) × Battery Capacity × Round-trip Efficiency

Time-of-Use Savings:
- Charge during off-peak: $0.08/kWh
- Discharge during peak: $0.32/kWh
- Daily savings: (0.32 - 0.08) × 10 kWh × 0.92 = $2.21

Annual Savings = Daily Savings × 365 = $807
```

### Backup Power Duration
```
Backup Duration (hours) = Battery Capacity (kWh) × DoD × Efficiency / Average Load (kW)

Critical Load Calculation:
- Refrigerator: 150W average
- Lights (LED): 100W
- Internet/Communications: 50W
- Well pump: 750W (intermittent)
- Total critical load: ~1 kW average

10 kWh battery provides: 10 × 0.90 × 0.95 / 1 = 8.5 hours
```

## 📈 Carbon Footprint Calculations

### Solar Carbon Offset
```
Annual CO₂ Reduction (lbs) = Annual Production (kWh) × Grid Emission Factor (lbs CO₂/kWh)

US Grid Average: 0.92 lbs CO₂/kWh
Coal-heavy regions: 1.5-2.0 lbs CO₂/kWh
Clean grids: 0.2-0.5 lbs CO₂/kWh

Example:
10 kW system × 1,400 kWh/kW/year × 0.92 lbs/kWh = 12,880 lbs CO₂/year
Equivalent to planting 160 trees annually
```

### Heat Pump Carbon Savings
```
Gas Furnace Emissions = (Therms × 11.7 lbs CO₂/therm)
Heat Pump Emissions = (kWh × Grid Factor)

Example Comparison:
- Gas furnace (95% eff): 1,000 therms × 11.7 = 11,700 lbs CO₂
- Heat pump (COP 3.0): 3,400 kWh × 0.92 = 3,128 lbs CO₂
- Reduction: 8,572 lbs CO₂/year (73%)
```

## 💰 Incentive Calculations

### Federal Tax Credit (30%)
```
Tax Credit = System Cost × 0.30

Eligible Costs:
- Solar panels and inverters
- Battery storage (if charged by solar)
- Installation labor
- Permits and inspections
- Sales tax on equipment

Example:
$25,000 system × 0.30 = $7,500 tax credit
```

### Utility Rebates
```
Common Rebate Structures:
- Fixed $/kW installed (e.g., $500/kW)
- Performance-based ($/kWh produced)
- Time-limited offers

Net Cost = Gross Cost - Federal Credit - State Rebates - Utility Incentives
```

### Depreciation (Commercial)
```
MACRS 5-Year Schedule:
Year 1: 20% + 80% bonus = 100% (current law)
Or standard: 20%, 32%, 19.2%, 11.52%, 11.52%, 5.76%

Tax Savings = Depreciation × Tax Rate
```

## 🏗️ Quick Reference Tables

### Solar Production by Region (kWh/kW/year)
```
Southwest (AZ, NV): 1,800-2,000
California: 1,650-1,850
Southeast: 1,400-1,550
Northeast: 1,200-1,350
Pacific Northwest: 1,100-1,250
```

### Electricity Rate Escalation
```
Historical average: 2.5-3.5%/year
Recent trends: 4-7%/year
Conservative estimate: 3%/year
Aggressive estimate: 5%/year
```

### Equipment Lifespan
```
Solar panels: 25-30 years (0.5% degradation/year)
Inverters: 10-15 years (may need 1 replacement)
Batteries: 10-15 years (capacity based)
Heat pumps: 15-20 years
```

## 📝 Customer Presentation Templates

### Simple Payback Statement
```
"With a net cost of $[X] after incentives, and annual savings of $[Y], 
your system will pay for itself in [Z] years. After that, you'll enjoy 
free electricity for the remaining 20+ years of the system's life."
```

### Monthly Cash Flow
```
"Your current electric bill: $[A]/month
Estimated solar loan payment: $[B]/month
Net monthly savings: $[A-B]/month from day one
After loan payoff: $[A]/month in pure savings"
```

### Environmental Impact
```
"Your [X] kW solar system will offset [Y] lbs of CO₂ annually, 
equivalent to:
- Planting [Y/80] trees every year
- Taking [Y/10,000] cars off the road
- Saving [Y/20] barrels of oil"
```

## 🔧 Advanced Calculations

### Load Shifting Value
```
Peak Demand Reduction = Battery Power (kW) × Demand Charge ($/kW)
Monthly Savings = Peak Reduction × 12

Example:
10 kW battery × $15/kW demand charge = $150/month = $1,800/year
```

### Virtual Power Plant Revenue
```
VPP Payment = (Discharge Events × kWh Discharged × $/kWh)

Typical programs:
- 20-50 events/year
- 2-4 hours/event
- $0.20-0.50/kWh payment
- Annual revenue: $400-1,000 for 10 kWh battery
```

### Microgrid Economics
```
Value Streams:
1. Energy arbitrage
2. Demand charge reduction
3. Backup power value
4. Grid services revenue
5. Carbon credits (where applicable)

Total Annual Value = Sum of all value streams
ROI = Total Annual Value / System Cost
```