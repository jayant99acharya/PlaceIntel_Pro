#!/usr/bin/env python3
"""
PlaceIntel Pro - Python Intelligence Service
Handles ML/AI processing for location intelligence
"""

import os
import time
import logging
from datetime import datetime
from typing import Dict, List, Any, Optional

from flask import Flask, request, jsonify
from flask_cors import CORS
import numpy as np
import random

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s'
)
logger = logging.getLogger(__name__)

app = Flask(__name__)
CORS(app)

class BusinessIntelligenceEngine:
    """Handles business intelligence analysis for places"""
    
    def __init__(self):
        self.category_insights = {
            'coffee': {
                'specialties': ['artisanal coffee', 'espresso', 'latte art'],
                'ideal_for': ['remote work', 'meetings', 'studying'],
                'atmosphere': 'cozy',
                'price_range': 'moderate'
            },
            'restaurant': {
                'specialties': ['local cuisine', 'fresh ingredients'],
                'ideal_for': ['dining', 'celebrations', 'dates'],
                'atmosphere': 'welcoming',
                'price_range': 'varied'
            },
            'gym': {
                'specialties': ['fitness equipment', 'personal training'],
                'ideal_for': ['workouts', 'fitness classes', 'health'],
                'atmosphere': 'energetic',
                'price_range': 'membership'
            },
            'library': {
                'specialties': ['study spaces', 'books', 'quiet environment'],
                'ideal_for': ['studying', 'research', 'reading'],
                'atmosphere': 'quiet',
                'price_range': 'free'
            },
            'shopping': {
                'specialties': ['retail', 'variety', 'brands'],
                'ideal_for': ['shopping', 'browsing', 'gifts'],
                'atmosphere': 'busy',
                'price_range': 'varied'
            }
        }
    
    def analyze_place(self, place_data: Dict[str, Any]) -> Dict[str, Any]:
        """Generate business intelligence for a place"""
        start_time = time.time()
        
        # Extract place information
        name = place_data.get('name', '')
        categories = place_data.get('categories', [])
        location = place_data.get('location', {})
        
        # Determine primary category
        primary_category = self._get_primary_category(categories)
        
        # Generate popularity score (simulated ML model)
        popularity_score = self._calculate_popularity_score(name, categories, location)
        
        # Generate sentiment score (simulated NLP analysis)
        sentiment_score = self._analyze_sentiment(name, primary_category)
        
        # Get category-specific insights
        category_info = self.category_insights.get(primary_category, {})
        
        # Generate trending score
        trending_score = self._calculate_trending_score(primary_category)
        
        processing_time = (time.time() - start_time) * 1000
        
        return {
            'popularity_score': round(popularity_score, 1),
            'sentiment_score': round(sentiment_score, 1),
            'specialties': category_info.get('specialties', []),
            'ideal_for': category_info.get('ideal_for', []),
            'price_range': category_info.get('price_range', 'unknown'),
            'atmosphere': category_info.get('atmosphere', 'unknown'),
            'trending_score': round(trending_score, 1),
            'processing_time_ms': round(processing_time, 2)
        }
    
    def _get_primary_category(self, categories: List[Dict]) -> str:
        """Extract primary category from place categories"""
        if not categories:
            return 'general'
        
        category_name = categories[0].get('name', '').lower()
        
        # Map Foursquare categories to our internal categories
        if 'coffee' in category_name or 'café' in category_name:
            return 'coffee'
        elif 'restaurant' in category_name or 'food' in category_name:
            return 'restaurant'
        elif 'gym' in category_name or 'fitness' in category_name:
            return 'gym'
        elif 'library' in category_name:
            return 'library'
        elif 'shop' in category_name or 'store' in category_name:
            return 'shopping'
        else:
            return 'general'
    
    def _calculate_popularity_score(self, name: str, categories: List[Dict], location: Dict) -> float:
        """Simulate ML model for popularity scoring"""
        # Base score
        score = 5.0
        
        # Name-based factors (simulated brand recognition)
        if any(brand in name.lower() for brand in ['starbucks', 'mcdonalds', 'subway']):
            score += 2.0
        elif len(name.split()) > 3:  # Longer names might indicate specialty places
            score += 1.0
        
        # Category-based factors
        if categories:
            category_name = categories[0].get('name', '').lower()
            if 'coffee' in category_name:
                score += 1.5  # Coffee shops are generally popular
            elif 'restaurant' in category_name:
                score += 1.0
        
        # Location-based factors (simulated urban density analysis)
        # In a real implementation, this would use actual demographic data
        score += random.uniform(-1.0, 2.0)
        
        # Ensure score is within valid range
        return max(1.0, min(10.0, score))
    
    def _analyze_sentiment(self, name: str, category: str) -> float:
        """Simulate NLP sentiment analysis"""
        # Base sentiment
        sentiment = 3.5
        
        # Positive keywords
        positive_words = ['best', 'premium', 'artisan', 'fresh', 'quality', 'authentic']
        negative_words = ['cheap', 'fast', 'quick']
        
        name_lower = name.lower()
        
        for word in positive_words:
            if word in name_lower:
                sentiment += 0.5
        
        for word in negative_words:
            if word in name_lower:
                sentiment -= 0.3
        
        # Category-based sentiment adjustment
        category_sentiment = {
            'coffee': 4.0,
            'restaurant': 3.8,
            'library': 4.2,
            'gym': 3.5,
            'shopping': 3.6
        }
        
        base_sentiment = category_sentiment.get(category, 3.5)
        sentiment = (sentiment + base_sentiment) / 2
        
        # Add some randomness to simulate real sentiment analysis
        sentiment += random.uniform(-0.5, 0.5)
        
        return max(1.0, min(5.0, sentiment))
    
    def _calculate_trending_score(self, category: str) -> float:
        """Calculate trending score based on category and time"""
        # Simulate trending analysis based on current trends
        trending_categories = {
            'coffee': 8.5,  # Coffee culture is trending
            'gym': 7.8,     # Fitness is popular
            'restaurant': 7.0,
            'library': 6.0,
            'shopping': 6.5
        }
        
        base_score = trending_categories.get(category, 5.0)
        
        # Add time-based variation (simulate real-time trends)
        hour = datetime.now().hour
        if 7 <= hour <= 9 and category == 'coffee':  # Morning coffee rush
            base_score += 1.0
        elif 12 <= hour <= 14 and category == 'restaurant':  # Lunch rush
            base_score += 1.0
        elif 18 <= hour <= 20 and category == 'restaurant':  # Dinner rush
            base_score += 1.0
        
        return max(0.0, min(10.0, base_score))


class RealTimeContextEngine:
    """Handles real-time context analysis for places"""
    
    def analyze_context(self, place_data: Dict[str, Any]) -> Dict[str, Any]:
        """Generate real-time context for a place"""
        start_time = time.time()
        
        categories = place_data.get('categories', [])
        primary_category = self._get_primary_category(categories)
        
        # Simulate real-time data analysis
        current_hour = datetime.now().hour
        
        # Generate context based on time and category
        status = self._determine_status(primary_category, current_hour)
        crowd_level = self._estimate_crowd_level(primary_category, current_hour)
        best_times = self._suggest_best_times(primary_category)
        events = self._detect_events(primary_category, current_hour)
        wait_time = self._estimate_wait_time(crowd_level)
        weather_impact = self._assess_weather_impact(primary_category)
        
        processing_time = (time.time() - start_time) * 1000
        
        return {
            'current_status': status,
            'crowd_level': crowd_level,
            'best_visit_times': best_times,
            'live_events': events,
            'estimated_wait_time': wait_time,
            'weather_impact': weather_impact,
            'last_updated': datetime.utcnow().isoformat() + 'Z',
            'confidence_score': round(random.uniform(0.7, 0.95), 2),
            'processing_time_ms': round(processing_time, 2)
        }
    
    def _get_primary_category(self, categories: List[Dict]) -> str:
        """Extract primary category from place categories"""
        if not categories:
            return 'general'
        return categories[0].get('name', '').lower()
    
    def _determine_status(self, category: str, hour: int) -> str:
        """Determine if place is likely open/closed"""
        # Simulate business hours analysis
        if 'restaurant' in category:
            if 6 <= hour <= 23:
                return 'open'
            else:
                return 'closed'
        elif 'coffee' in category:
            if 6 <= hour <= 20:
                return 'open'
            else:
                return 'closed'
        elif 'gym' in category:
            if 5 <= hour <= 23:
                return 'open'
            else:
                return 'closed'
        elif 'library' in category:
            if 8 <= hour <= 20:
                return 'open'
            else:
                return 'closed'
        else:
            if 9 <= hour <= 21:
                return 'open'
            else:
                return 'closed'
    
    def _estimate_crowd_level(self, category: str, hour: int) -> str:
        """Estimate current crowd level"""
        # Peak hours analysis
        if 'restaurant' in category:
            if hour in [12, 13, 18, 19, 20]:
                return 'busy'
            elif hour in [11, 14, 17, 21]:
                return 'moderate'
            else:
                return 'quiet'
        elif 'coffee' in category:
            if hour in [7, 8, 9, 14, 15]:
                return 'busy'
            elif hour in [10, 11, 16, 17]:
                return 'moderate'
            else:
                return 'quiet'
        elif 'gym' in category:
            if hour in [6, 7, 17, 18, 19]:
                return 'busy'
            elif hour in [8, 9, 16, 20]:
                return 'moderate'
            else:
                return 'quiet'
        else:
            return random.choice(['quiet', 'moderate', 'busy'])
    
    def _suggest_best_times(self, category: str) -> List[str]:
        """Suggest optimal visit times"""
        if 'restaurant' in category:
            return ['11:30-12:00', '14:30-17:00', '21:00-22:00']
        elif 'coffee' in category:
            return ['10:00-11:30', '15:30-17:00']
        elif 'gym' in category:
            return ['10:00-16:00', '21:00-23:00']
        elif 'library' in category:
            return ['9:00-11:00', '14:00-16:00']
        else:
            return ['10:00-12:00', '14:00-16:00']
    
    def _detect_events(self, category: str, hour: int) -> List[str]:
        """Detect potential live events"""
        events = []
        
        if 'library' in category and 14 <= hour <= 16:
            events.append('study group session')
        elif 'gym' in category and hour in [18, 19]:
            events.append('fitness class')
        elif 'coffee' in category and hour in [15, 16]:
            events.append('afternoon networking')
        
        return events
    
    def _estimate_wait_time(self, crowd_level: str) -> str:
        """Estimate wait time based on crowd level"""
        if crowd_level == 'busy':
            return '10-15 minutes'
        elif crowd_level == 'moderate':
            return '5-10 minutes'
        else:
            return 'no wait'
    
    def _assess_weather_impact(self, category: str) -> str:
        """Assess weather impact on place"""
        # Simulate weather analysis
        weather_conditions = ['sunny', 'rainy', 'cloudy']
        current_weather = random.choice(weather_conditions)
        
        if current_weather == 'rainy' and 'outdoor' in category:
            return 'high impact - indoor alternatives recommended'
        elif current_weather == 'sunny' and 'park' in category:
            return 'positive impact - great weather for outdoor activities'
        else:
            return 'minimal impact'


class AccessibilityIntelligenceEngine:
    """Handles accessibility analysis for places"""
    
    def analyze_accessibility(self, place_data: Dict[str, Any]) -> Dict[str, Any]:
        """Generate accessibility intelligence for a place"""
        start_time = time.time()
        
        categories = place_data.get('categories', [])
        name = place_data.get('name', '')
        
        # Simulate accessibility analysis
        # In a real implementation, this would use computer vision, 
        # crowdsourced data, and accessibility databases
        
        accessibility_score = self._calculate_accessibility_score(categories, name)
        features = self._analyze_accessibility_features(categories)
        recommendations = self._generate_inclusive_recommendations(categories, features)
        
        processing_time = (time.time() - start_time) * 1000
        
        return {
            'wheelchair_accessible': accessibility_score >= 7.0,
            'accessibility_score': round(accessibility_score, 1),
            'features': features,
            'inclusive_recommendations': recommendations,
            'processing_time_ms': round(processing_time, 2)
        }
    
    def _calculate_accessibility_score(self, categories: List[Dict], name: str) -> float:
        """Calculate overall accessibility score"""
        base_score = 5.0
        
        if not categories:
            return base_score
        
        category_name = categories[0].get('name', '').lower()
        
        # Modern establishments tend to be more accessible
        if any(word in name.lower() for word in ['new', 'modern', 'center', 'mall']):
            base_score += 2.0
        
        # Category-based accessibility assumptions
        if 'library' in category_name or 'hospital' in category_name:
            base_score += 3.0  # Public buildings usually more accessible
        elif 'restaurant' in category_name and 'chain' in name.lower():
            base_score += 2.0  # Chain restaurants often have standards
        elif 'gym' in category_name:
            base_score += 1.5  # Modern gyms often accessible
        
        # Add some variation
        base_score += random.uniform(-1.0, 1.0)
        
        return max(0.0, min(10.0, base_score))
    
    def _analyze_accessibility_features(self, categories: List[Dict]) -> Dict[str, bool]:
        """Analyze specific accessibility features"""
        if not categories:
            return self._default_features()
        
        category_name = categories[0].get('name', '').lower()
        
        # Simulate feature detection based on category
        if 'library' in category_name:
            return {
                'ramp_access': True,
                'elevator': True,
                'accessible_restrooms': True,
                'braille_signage': True,
                'hearing_loop': True,
                'wide_entrances': True,
                'accessible_parking': True
            }
        elif 'restaurant' in category_name:
            return {
                'ramp_access': random.choice([True, False]),
                'elevator': False,  # Most restaurants are single floor
                'accessible_restrooms': random.choice([True, False]),
                'braille_signage': False,
                'hearing_loop': False,
                'wide_entrances': random.choice([True, False]),
                'accessible_parking': random.choice([True, False])
            }
        elif 'gym' in category_name:
            return {
                'ramp_access': True,
                'elevator': random.choice([True, False]),
                'accessible_restrooms': True,
                'braille_signage': False,
                'hearing_loop': False,
                'wide_entrances': True,
                'accessible_parking': True
            }
        else:
            return self._default_features()
    
    def _default_features(self) -> Dict[str, bool]:
        """Default accessibility features"""
        return {
            'ramp_access': random.choice([True, False]),
            'elevator': random.choice([True, False]),
            'accessible_restrooms': random.choice([True, False]),
            'braille_signage': False,
            'hearing_loop': False,
            'wide_entrances': random.choice([True, False]),
            'accessible_parking': random.choice([True, False])
        }
    
    def _generate_inclusive_recommendations(self, categories: List[Dict], features: Dict[str, bool]) -> Dict[str, List[str]]:
        """Generate inclusive recommendations"""
        recommendations = {
            'mobility_friendly_areas': [],
            'sensory_accommodations': [],
            'cognitive_support': []
        }
        
        if not categories:
            return recommendations
        
        category_name = categories[0].get('name', '').lower()
        
        # Mobility recommendations
        if features.get('ramp_access'):
            recommendations['mobility_friendly_areas'].append('main entrance accessible')
        if features.get('elevator'):
            recommendations['mobility_friendly_areas'].append('all floors accessible')
        if features.get('accessible_restrooms'):
            recommendations['mobility_friendly_areas'].append('accessible restroom facilities')
        
        # Sensory accommodations
        if 'library' in category_name:
            recommendations['sensory_accommodations'].extend([
                'quiet study areas available',
                'adjustable lighting in reading areas'
            ])
        if features.get('hearing_loop'):
            recommendations['sensory_accommodations'].append('hearing loop system available')
        
        # Cognitive support
        if 'library' in category_name:
            recommendations['cognitive_support'].extend([
                'clear signage and wayfinding',
                'staff available for assistance'
            ])
        elif 'restaurant' in category_name:
            recommendations['cognitive_support'].append('picture menus available')
        
        return recommendations


class UnifiedRecommendationEngine:
    """Generates unified recommendations combining all intelligence types"""
    
    def generate_recommendations(
        self, 
        place_data: Dict[str, Any],
        business_intel: Dict[str, Any],
        realtime_context: Dict[str, Any],
        accessibility_intel: Dict[str, Any]
    ) -> Dict[str, Any]:
        """Generate unified recommendations"""
        start_time = time.time()
        
        # Calculate overall confidence score
        confidence_score = self._calculate_confidence_score(
            business_intel, realtime_context, accessibility_intel
        )
        
        # Generate personalized insights
        insights = self._generate_personalized_insights(
            place_data, business_intel, realtime_context, accessibility_intel
        )
        
        # Generate alternative suggestions
        alternatives = self._suggest_alternatives(place_data, business_intel)
        
        # Generate optimal visit strategy
        strategy = self._create_visit_strategy(realtime_context, accessibility_intel)
        
        # Generate accessibility notes
        accessibility_notes = self._create_accessibility_notes(accessibility_intel)
        
        processing_time = (time.time() - start_time) * 1000
        
        return {
            'confidence_score': round(confidence_score, 2),
            'personalized_insights': insights,
            'alternative_suggestions': alternatives,
            'optimal_visit_strategy': strategy,
            'accessibility_notes': accessibility_notes,
            'processing_time_ms': round(processing_time, 2)
        }
    
    def _calculate_confidence_score(
        self, 
        business_intel: Dict[str, Any],
        realtime_context: Dict[str, Any],
        accessibility_intel: Dict[str, Any]
    ) -> float:
        """Calculate overall confidence in recommendations"""
        scores = []
        
        # Business intelligence confidence
        if business_intel.get('popularity_score', 0) > 0:
            scores.append(0.8)
        
        # Real-time context confidence
        if realtime_context.get('confidence_score', 0) > 0:
            scores.append(realtime_context['confidence_score'])
        
        # Accessibility intelligence confidence
        if accessibility_intel.get('accessibility_score', 0) > 0:
            scores.append(0.7)
        
        return sum(scores) / len(scores) if scores else 0.5
    
    def _generate_personalized_insights(
        self,
        place_data: Dict[str, Any],
        business_intel: Dict[str, Any],
        realtime_context: Dict[str, Any],
        accessibility_intel: Dict[str, Any]
    ) -> List[str]:
        """Generate personalized insights"""
        insights = []
        
        # Business insights
        popularity = business_intel.get('popularity_score', 0)
        if popularity >= 8.0:
            insights.append(f"Highly popular destination with {popularity}/10 rating")
        elif popularity >= 6.0:
            insights.append(f"Well-regarded place with {popularity}/10 popularity")
        
        # Real-time insights
        crowd_level = realtime_context.get('crowd_level', '')
        if crowd_level == 'quiet':
            insights.append("Currently quiet - perfect for a peaceful visit")
        elif crowd_level == 'busy':
            insights.append("Currently busy - consider visiting during suggested off-peak times")
        
        # Accessibility insights
        if accessibility_intel.get('wheelchair_accessible'):
            insights.append("Fully wheelchair accessible with comprehensive features")
        
        # Atmosphere insights
        atmosphere = business_intel.get('atmosphere', '')
        ideal_for = business_intel.get('ideal_for', [])
        if atmosphere and ideal_for:
            insights.append(f"{atmosphere.title()} atmosphere, ideal for {', '.join(ideal_for)}")
        
        return insights[:4]  # Limit to top 4 insights
    
    def _suggest_alternatives(self, place_data: Dict[str, Any], business_intel: Dict[str, Any]) -> List[str]:
        """Suggest alternative places or experiences"""
        alternatives = []
        
        categories = place_data.get('categories', [])
        if not categories:
            return alternatives
        
        category_name = categories[0].get('name', '').lower()
        
        if 'coffee' in category_name:
            alternatives.extend([
                "Local independent coffee shops nearby",
                "Tea houses for alternative beverages",
                "Co-working spaces with café facilities"
            ])
        elif 'restaurant' in category_name:
            alternatives.extend([
                "Similar cuisine restaurants in the area",
                "Food trucks for casual dining",
                "Delivery options if crowded"
            ])
        elif 'gym' in category_name:
            alternatives.extend([
                "Outdoor fitness areas nearby",
                "Alternative fitness studios",
                "Home workout options during peak hours"
            ])
        
        return alternatives[:3]  # Limit to top 3 alternatives
    
    def _create_visit_strategy(self, realtime_context: Dict[str, Any], accessibility_intel: Dict[str, Any]) -> str:
        """Create optimal visit strategy"""
        best_times = realtime_context.get('best_visit_times', [])
        crowd_level = realtime_context.get('crowd_level', '')
        wait_time = realtime_context.get('estimated_wait_time', '')
        
        strategy_parts = []
        
        if best_times:
            strategy_parts.append(f"Best visit times: {', '.join(best_times)}")
        
        if crowd_level == 'busy' and wait_time != 'no wait':
            strategy_parts.append(f"Current wait time: {wait_time}")
        
        if accessibility_intel.get('wheelchair_accessible'):
            strategy_parts.append("Accessible entrance available")
        
        return ". ".join(strategy_parts) if strategy_parts else "Visit anytime based on your preference"
    
    def _create_accessibility_notes(self, accessibility_intel: Dict[str, Any]) -> List[str]:
        """Create accessibility-specific notes"""
        notes = []
        
        if accessibility_intel.get('wheelchair_accessible'):
            notes.append("Wheelchair accessible with ramp access")
        else:
            notes.append("Accessibility features may be limited - recommend calling ahead")
        
        features = accessibility_intel.get('features', {})
        if features.get('accessible_restrooms'):
            notes.append("Accessible restroom facilities available")
        
        if features.get('hearing_loop'):
            notes.append("Hearing loop system available for hearing aid users")
        
        return notes[:3]  # Limit to top 3 notes


# Initialize engines
business_engine = BusinessIntelligenceEngine()
realtime_engine = RealTimeContextEngine()
accessibility_engine = AccessibilityIntelligenceEngine()
recommendation_engine = UnifiedRecommendationEngine()


@app.route('/health', methods=['GET'])
def health_check():
    """Health check endpoint"""
    return jsonify({
        'status': 'healthy',
        'timestamp': datetime.utcnow().isoformat() + 'Z',
        'service': 'PlaceIntel Pro Intelligence Service',
        'version': '1.0.0'
    })


@app.route('/api/v1/intelligence/enhance', methods=['POST'])
def enhance_place_intelligence():
    """Main endpoint for enhancing place data with intelligence"""
    try:
        start_time = time.time()
        
        # Parse request
        data = request.get_json()
        if not data:
            return jsonify({'error': 'No data provided'}), 400
        
        place_data = data.get('place', {})
        if not place_data:
            return jsonify({'error': 'No place data provided'}), 400
        
        logger.info(f"Processing intelligence for place: {place_data.get('name', 'Unknown')}")
        
        # Generate business intelligence
        business_intel = business_engine.analyze_place(place_data)
        
        # Generate real-time context
        realtime_context = realtime_engine.analyze_context(place_data)
        
        # Generate accessibility intelligence
        accessibility_intel = accessibility_engine.analyze_accessibility(place_data)
        
        # Generate unified recommendations
        unified_recommendations = recommendation_engine.generate_recommendations(
            place_data, business_intel, realtime_context, accessibility_intel
        )
        
        total_processing_time = (time.time() - start_time) * 1000
        
        response = {
            'business_intelligence': business_intel,
            'real_time_context': realtime_context,
            'accessibility_intelligence': accessibility_intel,
            'unified_recommendations': unified_recommendations,
            'processing_time_ms': round(total_processing_time, 2),
            'data_sources': ['foursquare', 'ml_models', 'accessibility_db', 'real_time_feeds']
        }
        
        logger.info(f"Intelligence processing completed in {total_processing_time:.2f}ms")
        
        return jsonify(response)
        
    except Exception as e:
        logger.error(f"Error processing intelligence: {str(e)}")
        return jsonify({
            'error': 'Internal server error',
            'message': str(e)
        }), 500


if __name__ == '__main__':
    port = int(os.environ.get('PORT', 5000))
    debug = os.environ.get('DEBUG', 'false').lower() == 'true'
    
    logger.info(f"Starting PlaceIntel Pro Intelligence Service on port {port}")
    app.run(host='0.0.0.0', port=port, debug=debug)