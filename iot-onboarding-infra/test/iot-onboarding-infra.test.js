"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const assert_1 = require("@aws-cdk/assert");
const cdk = require("@aws-cdk/core");
const IotOnboardingInfra = require("../lib/iot-onboarding-infra-stack");
test('Empty Stack', () => {
    const app = new cdk.App();
    // WHEN
    const stack = new IotOnboardingInfra.IOTOnboardingInfraStack(app, 'MyTestStack');
    // THEN
    (0, assert_1.expect)(stack).to((0, assert_1.matchTemplate)({
        "Resources": {}
    }, assert_1.MatchStyle.EXACT));
});
//# sourceMappingURL=data:application/json;base64,eyJ2ZXJzaW9uIjozLCJmaWxlIjoiaW90LW9uYm9hcmRpbmctaW5mcmEudGVzdC5qcyIsInNvdXJjZVJvb3QiOiIiLCJzb3VyY2VzIjpbImlvdC1vbmJvYXJkaW5nLWluZnJhLnRlc3QudHMiXSwibmFtZXMiOltdLCJtYXBwaW5ncyI6Ijs7QUFBQSw0Q0FBaUY7QUFDakYscUNBQXFDO0FBQ3JDLHdFQUF3RTtBQUV4RSxJQUFJLENBQUMsYUFBYSxFQUFFLEdBQUcsRUFBRTtJQUN2QixNQUFNLEdBQUcsR0FBRyxJQUFJLEdBQUcsQ0FBQyxHQUFHLEVBQUUsQ0FBQztJQUMxQixPQUFPO0lBQ1AsTUFBTSxLQUFLLEdBQUcsSUFBSSxrQkFBa0IsQ0FBQyx1QkFBdUIsQ0FBQyxHQUFHLEVBQUUsYUFBYSxDQUFDLENBQUM7SUFDakYsT0FBTztJQUNQLElBQUEsZUFBUyxFQUFDLEtBQUssQ0FBQyxDQUFDLEVBQUUsQ0FBQyxJQUFBLHNCQUFhLEVBQUM7UUFDaEMsV0FBVyxFQUFFLEVBQUU7S0FDaEIsRUFBRSxtQkFBVSxDQUFDLEtBQUssQ0FBQyxDQUFDLENBQUE7QUFDdkIsQ0FBQyxDQUFDLENBQUMiLCJzb3VyY2VzQ29udGVudCI6WyJpbXBvcnQgeyBleHBlY3QgYXMgZXhwZWN0Q0RLLCBtYXRjaFRlbXBsYXRlLCBNYXRjaFN0eWxlIH0gZnJvbSAnQGF3cy1jZGsvYXNzZXJ0JztcbmltcG9ydCAqIGFzIGNkayBmcm9tICdAYXdzLWNkay9jb3JlJztcbmltcG9ydCAqIGFzIElvdE9uYm9hcmRpbmdJbmZyYSBmcm9tICcuLi9saWIvaW90LW9uYm9hcmRpbmctaW5mcmEtc3RhY2snO1xuXG50ZXN0KCdFbXB0eSBTdGFjaycsICgpID0+IHtcbiAgY29uc3QgYXBwID0gbmV3IGNkay5BcHAoKTtcbiAgLy8gV0hFTlxuICBjb25zdCBzdGFjayA9IG5ldyBJb3RPbmJvYXJkaW5nSW5mcmEuSU9UT25ib2FyZGluZ0luZnJhU3RhY2soYXBwLCAnTXlUZXN0U3RhY2snKTtcbiAgLy8gVEhFTlxuICBleHBlY3RDREsoc3RhY2spLnRvKG1hdGNoVGVtcGxhdGUoe1xuICAgIFwiUmVzb3VyY2VzXCI6IHt9XG4gIH0sIE1hdGNoU3R5bGUuRVhBQ1QpKVxufSk7XG4iXX0=