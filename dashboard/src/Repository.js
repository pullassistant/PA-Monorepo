import React, {Component} from 'react';
import './Repository.css';
import image from './img/repository.png';

class Repository extends Component {

    render() {
        return (
            <div className="repository">
                <img src={image} alt="Repository"/> <a href={this.props.repository.html_url}
                                                       target="_blank">{this.props.repository.name}</a>
            </div>
        );
    }
}

export default Repository;
